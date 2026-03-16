import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.enableCors({
    origin: [
      'http://localhost:3000',
      'https://www.climia.com.br',
      'https://climia.com.br',
    ],
  });

  const port = process.env.PORT ?? 3001;
  await app.listen(port);
  console.log(`ClimIA API running on http://localhost:${port}`);
}

bootstrap();
