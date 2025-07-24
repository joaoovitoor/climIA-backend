#!/bin/bash

echo "🔍 MONITORAMENTO DE CUSTOS E RECURSOS AWS"
echo "=========================================="

echo ""
echo "📊 LAMBDA FUNCTION:"
aws lambda get-function --function-name climia-api --query 'Configuration.{FunctionName:FunctionName,Runtime:Runtime,MemorySize:MemorySize,Timeout:Timeout,CodeSize:CodeSize}' --output table

echo ""
echo "🌐 API GATEWAY:"
aws apigateway get-rest-api --rest-api-id hls852t472 --query '{Name:name,Description:description,CreatedDate:createdDate,Version:version}' --output table

echo ""
echo "📝 LOGS:"
aws logs describe-log-groups --log-group-name-prefix "/aws/lambda/climia-api" --query 'logGroups[0].{LogGroupName:logGroupName,StoredBytes:storedBytes,RetentionInDays:retentionInDays}' --output table

echo ""
echo "💰 ESTIMATIVA DE CUSTOS:"
echo "Lambda: ~$0.20 por 1M de invocações (512MB)"
echo "API Gateway: ~$3.50 por 1M de requests"
echo "CloudWatch Logs: ~$0.50 por GB armazenado"
echo ""
echo "📈 PARA MONITORAR CUSTOS EM TEMPO REAL:"
echo "1. Acesse: https://console.aws.amazon.com/cost-management/"
echo "2. Vá em 'Cost Explorer'"
echo "3. Filtre por serviço: Lambda, API Gateway, CloudWatch"
echo ""
echo "🔔 PARA CRIAR ALERTAS DE CUSTO:"
echo "1. Acesse: https://console.aws.amazon.com/cloudwatch/"
echo "2. Vá em 'Alarms' > 'Create alarm'"
echo "3. Configure alertas para:"
echo "   - Lambda Duration > 10s"
echo "   - Lambda Errors > 0"
echo "   - API Gateway 4XX/5XX errors"
echo ""
echo "📋 COMANDOS ÚTEIS:"
echo "aws lambda get-function --function-name climia-api"
echo "aws logs describe-log-streams --log-group-name /aws/lambda/climia-api"
echo "aws apigateway get-stages --rest-api-id hls852t472" 