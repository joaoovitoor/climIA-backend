import {
  Controller,
  Get,
  Query,
  BadRequestException,
  InternalServerErrorException,
  Headers,
  UnauthorizedException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { WeatherService, Forecast } from './weather.service';

@Controller()
export class WeatherController {
  constructor(
    private readonly weatherService: WeatherService,
    private readonly configService: ConfigService,
  ) {}

  @Get('health')
  health() {
    return { status: 'ok', timestamp: new Date().toISOString() };
  }

  @Get()
  async getForecast(
    @Headers('authorization') authHeader: string,
    @Query('cidade') cidade: string,
    @Query('estado') estado: string,
    @Query('data') data: string,
    @Query('datainicio') dataInicio: string,
    @Query('datafim') dataFim: string,
  ): Promise<Forecast[]> {
    this.validateAuth(authHeader);

    if (!cidade) {
      throw new BadRequestException('Parâmetro "cidade" é obrigatório');
    }

    const startDate = data || dataInicio;
    const endDate = data || dataFim;

    if (!startDate) {
      throw new BadRequestException(
        'Informe "data" ou "datainicio" e "datafim"',
      );
    }

    const cityData = this.findCity(cidade, estado);
    if (!cityData) {
      throw new BadRequestException(`Cidade "${cidade}" não encontrada`);
    }

    try {
      return await this.weatherService.getForecast(
        cityData.nome,
        cityData.uf,
        cityData.latitude,
        cityData.longitude,
        startDate,
        endDate || startDate,
      );
    } catch (error) {
      throw new InternalServerErrorException(
        'Erro ao buscar previsão: ' + (error as Error).message,
      );
    }
  }

  private validateAuth(authHeader: string) {
    const token = this.configService.get<string>('API_TOKEN');
    if (!token) return;

    if (!authHeader?.startsWith('Bearer ')) {
      throw new UnauthorizedException('Token de autenticação necessário');
    }

    if (authHeader.replace('Bearer ', '') !== token) {
      throw new UnauthorizedException('Token inválido');
    }
  }

  private citiesCache: Array<{
    nome: string;
    uf: string;
    latitude: number;
    longitude: number;
  }> | null = null;

  private findCity(cidade: string, estado?: string) {
    if (!this.citiesCache) {
      // eslint-disable-next-line @typescript-eslint/no-require-imports
      this.citiesCache = require('../../../frontend/data/cities.json');
    }

    const query = cidade
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '');

    return this.citiesCache!.find((c) => {
      const name = c.nome
        .toLowerCase()
        .normalize('NFD')
        .replace(/[\u0300-\u036f]/g, '');
      const matchName = name === query || name.includes(query);
      const matchUf =
        !estado || c.uf.toLowerCase() === estado.toLowerCase();
      return matchName && matchUf;
    });
  }
}
