import './trace'; // Must be first — starts OTEL SDK before http/express are loaded
import { ValidationPipe } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
const app = await NestFactory.create(AppModule);

app.setGlobalPrefix('api/v1');

app.useGlobalPipes(
new ValidationPipe({
whitelist: true,
transform: true,
}),
);

await app.listen(8001);

console.log('Point service is running on http://localhost:8001');
}

bootstrap();