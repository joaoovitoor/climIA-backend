import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, QueryCommand } from '@aws-sdk/lib-dynamodb';

export interface Forecast {
  cidade: string;
  estado: string;
  data: string;
  temperatura_minima: number;
  temperatura_maxima: number;
  temperatura_media: number;
  precipitacao: number;
}

interface CacheEntry {
  data: Forecast[];
  cachedAt: number;
}

@Injectable()
export class WeatherService implements OnModuleInit {
  private readonly logger = new Logger(WeatherService.name);
  private dynamo: DynamoDBDocumentClient;
  private readonly tableName: string;
  private readonly cache = new Map<string, CacheEntry>();
  private readonly CACHE_TTL = 24 * 60 * 60 * 1000;

  constructor(private readonly configService: ConfigService) {
    this.tableName =
      this.configService.get<string>('DYNAMODB_TABLE') ||
      this.configService.get<string>('DYNAMODB_TABLE_NAME') ||
      'ClimIA-Previsoes';
  }

  onModuleInit() {
    const client = new DynamoDBClient({
      region: this.configService.get<string>('AWS_REGION') || 'us-east-1',
    });
    this.dynamo = DynamoDBDocumentClient.from(client);
    this.logger.log(`DynamoDB conectado — tabela: ${this.tableName}`);
  }

  async getForecast(
    city: string,
    uf: string,
    latitude: number,
    longitude: number,
    startDate: string,
    endDate: string,
  ): Promise<Forecast[]> {
    const cacheKey = `${city}:${uf}:${startDate}:${endDate}`;
    const cached = this.cache.get(cacheKey);

    if (cached && Date.now() - cached.cachedAt < this.CACHE_TTL) {
      return cached.data;
    }

    const pk = `CITY#${city}#${uf}`;
    const dates = this.expandDateRange(startDate, endDate);
    const items = await this.queryAllDays(pk);

    const skSet = new Map(
      dates.map((d) => [`DATE#${d.day}#${d.month}`, d.dateStr]),
    );

    const forecasts: Forecast[] = [];

    for (const item of items) {
      const dateStr = skSet.get(item.SK);
      if (!dateStr) continue;

      forecasts.push({
        cidade: item.cidade,
        estado: item.estado,
        data: dateStr,
        temperatura_minima: Math.round(Number(item.temperatura_minima)),
        temperatura_maxima: Math.round(Number(item.temperatura_maxima)),
        temperatura_media: Math.round(Number(item.temperatura_media)),
        precipitacao: Number(item.precipitacao),
      });
    }

    forecasts.sort(
      (a, b) => new Date(a.data).getTime() - new Date(b.data).getTime(),
    );

    this.cache.set(cacheKey, { data: forecasts, cachedAt: Date.now() });

    return forecasts;
  }

  private cityCache = new Map<string, { items: any[]; cachedAt: number }>();

  private async queryAllDays(pk: string): Promise<any[]> {
    const cached = this.cityCache.get(pk);
    if (cached && Date.now() - cached.cachedAt < this.CACHE_TTL) {
      return cached.items;
    }

    const items: any[] = [];
    let lastKey: Record<string, any> | undefined;

    do {
      const result = await this.dynamo.send(
        new QueryCommand({
          TableName: this.tableName,
          KeyConditionExpression: 'PK = :pk',
          ExpressionAttributeValues: { ':pk': pk },
          ExclusiveStartKey: lastKey,
        }),
      );
      items.push(...(result.Items || []));
      lastKey = result.LastEvaluatedKey;
    } while (lastKey);

    this.cityCache.set(pk, { items, cachedAt: Date.now() });

    return items;
  }

  private expandDateRange(
    start: string,
    end: string,
  ): { day: string; month: string; dateStr: string }[] {
    const dates: { day: string; month: string; dateStr: string }[] = [];
    const current = new Date(start + 'T00:00:00');
    const endDate = new Date(end + 'T00:00:00');

    while (current <= endDate) {
      const day = String(current.getDate()).padStart(2, '0');
      const month = String(current.getMonth() + 1).padStart(2, '0');
      dates.push({
        day,
        month,
        dateStr: current.toISOString().split('T')[0],
      });
      current.setDate(current.getDate() + 1);
    }

    return dates;
  }
}
