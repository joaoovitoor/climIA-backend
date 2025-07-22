# 🔐 Configuração dos GitHub Secrets

## 📋 Secrets Necessários

Vá em: `Settings` → `Secrets and variables` → `Actions`

### **1. AWS Credentials**

```
AWS_ACCESS_KEY_ID=sua_access_key
AWS_SECRET_ACCESS_KEY=sua_secret_key
AWS_ACCOUNT_ID=123456789012
```

### **2. Lambda Role**

```
LAMBDA_ROLE_ARN=arn:aws:iam::123456789012:role/lambda-execution-role
```

### **3. Database Variables**

```
DB_HOST=seu-mysql-host.com
DB_PORT=3306
DB_USER=root
DB_PASS=sua-senha
DB_NAME=climia
```

## 🔧 Como Criar os Secrets

### **1. AWS Access Keys**

```bash
# Criar usuário IAM
aws iam create-user --user-name github-actions

# Criar access keys
aws iam create-access-key --user-name github-actions

# Anexar políticas
aws iam attach-user-policy \
  --user-name github-actions \
  --policy-arn arn:aws:iam::aws:policy/AWSLambdaFullAccess

aws iam attach-user-policy \
  --user-name github-actions \
  --policy-arn arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator
```

### **2. Lambda Role**

```bash
# Criar role
aws iam create-role \
  --role-name lambda-execution-role \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }'

# Anexar políticas
aws iam attach-role-policy \
  --role-name lambda-execution-role \
  --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

### **3. Account ID**

```bash
# Pegar Account ID
aws sts get-caller-identity --query Account --output text
```

## 🚀 Como Funciona o CI/CD

### **Fluxo Automático:**

1. **Push para main** → Trigger do workflow
2. **Test** → Roda testes
3. **Build** → Compila para Lambda
4. **Deploy** → Atualiza Lambda + API Gateway

### **Segurança:**

- ✅ **Secrets criptografados** - Não aparecem nos logs
- ✅ **Branch protection** - Só deploy na main
- ✅ **Testes obrigatórios** - Só deploy se testes passarem

## 📊 Monitoramento

### **Verificar Deploy:**

```bash
# Ver logs da Lambda
aws logs tail /aws/lambda/climia-api --follow

# Ver status
aws lambda get-function --function-name climia-api
```

### **Rollback:**

```bash
# Voltar versão anterior
aws lambda update-function-code \
  --function-name climia-api \
  --zip-file fileb://previous-version.zip
```

## 🎯 Próximos Passos

1. **Configurar secrets** no GitHub
2. **Fazer push** para main
3. **Verificar deploy** automático
4. **Testar endpoints** da API

## 🔍 Troubleshooting

### **Erro de Permissão:**

```bash
# Verificar se usuário tem permissões
aws sts get-caller-identity
```

### **Erro de Role:**

```bash
# Verificar se role existe
aws iam get-role --role-name lambda-execution-role
```

### **Erro de Database:**

- Verificar se MySQL está acessível
- Verificar variáveis de ambiente
- Verificar logs da Lambda
