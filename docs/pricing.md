# YOLO Pricing Guide

## Table of Contents
- [Overview](#overview)
- [Pricing Models](#pricing-models)
- [When to Use Your Own API Keys](#when-to-use-your-own-api-keys)
- [YOLO Packages](#yolo-packages)
- [License Management](#license-management)
- [Cost Comparison](#cost-comparison)
- [Enterprise Solutions](#enterprise-solutions)

## Overview

YOLO offers flexible pricing options to suit different development needs:

1. **Self-Hosted (BYO API Keys)**
   - Use your own OpenAI API keys
   - Full control over API usage
   - Pay directly to OpenAI

2. **YOLO Packages**
   - Pre-purchased token packages
   - Significant cost savings
   - Built-in license management

3. **Enterprise**
   - Custom solutions
   - Volume discounts
   - Priority support

## Pricing Models

### Self-Hosted Model
- You provide OpenAI API keys
- Pay OpenAI directly for API usage
- No additional YOLO costs
- Full control over model selection

### YOLO Packages
```
Starter:     $29/month  (100K tokens)
Developer:   $99/month  (500K tokens)
Team:        $299/month (2M tokens)
Business:    $999/month (10M tokens)
Enterprise:  Custom pricing
```

All packages include:
- Token pooling across team
- Rollover unused tokens (up to 3 months)
- Access to optimized prompts
- Basic support

## When to Use Your Own API Keys

### Use Your Own Keys When:

1. **Development/Testing**
   - Early project phases
   - Experimenting with features
   - Testing different models
   - Less than 50K tokens/month usage

2. **Specific Model Requirements**
   - Need custom OpenAI models
   - Using other AI providers
   - Special rate agreements with OpenAI

3. **Cost Considerations**
   ```
   OpenAI Direct Costs (GPT-4):
   - Input:  $0.03/1K tokens
   - Output: $0.06/1K tokens
   
   Average YOLO Command Costs:
   yolo commit:     ~1.5K tokens ($0.09)
   yolo feature:    ~2K tokens  ($0.12)
   yolo task:       ~1K tokens  ($0.06)
   ```

### Use YOLO Packages When:

1. **Production Usage**
   - Regular development workflow
   - Team of developers
   - Predictable usage patterns
   - More than 50K tokens/month

2. **Cost Optimization**
   - Need pooled resources
   - Want to cap spending
   - Benefit from bulk pricing

3. **Team Features**
   - Shared token pool
   - Usage monitoring
   - License management

## YOLO Packages

### Starter Package ($29/month)
- 100K tokens monthly
- Perfect for individual developers
- Basic support
- Features:
  - All core YOLO features
  - Single developer license
  - Email support

### Developer Package ($99/month)
- 500K tokens monthly
- Ideal for active developers
- Standard support
- Features:
  - Everything in Starter
  - Advanced prompts
  - Usage analytics
  - Priority email support

### Team Package ($299/month)
- 2M tokens monthly
- Great for small teams
- Premium support
- Features:
  - Everything in Developer
  - Up to 5 developer licenses
  - Team management
  - Slack support
  - Custom prompts

### Business Package ($999/month)
- 10M tokens monthly
- Enterprise-grade solution
- 24/7 support
- Features:
  - Everything in Team
  - Unlimited developers
  - Custom integrations
  - Dedicated support
  - Training sessions

### Enterprise (Custom)
- Custom token allocation
- Tailored solutions
- White-glove support
- Features:
  - Custom feature development
  - On-premises options
  - SLA guarantees
  - Account management

## License Management

### Activating Your License

1. **Via CLI**
   ```bash
   # Activate with license key
   yolo license activate XXXX-XXXX-XXXX-XXXX
   
   # View license status
   yolo license status
   
   # Deactivate license
   yolo license deactivate
   ```

2. **Via Configuration File**
   ```yaml
   # .yolo/config.yaml
   license:
     key: XXXX-XXXX-XXXX-XXXX
     type: team
     expires: 2025-12-31
   ```

### Managing Team Licenses

1. **Adding Team Members**
   ```bash
   yolo license add-user user@company.com
   yolo license list-users
   yolo license remove-user user@company.com
   ```

2. **Usage Monitoring**
   ```bash
   yolo usage report
   yolo usage by-user
   yolo usage by-command
   ```

## Cost Comparison

### Example Scenarios

1. **Solo Developer**
   ```
   Monthly Usage:
   - 100 commits (~150K tokens)
   - 20 features (~40K tokens)
   - 50 tasks (~50K tokens)
   
   OpenAI Direct Cost: ~$14.40
   YOLO Starter Package: $29 (better value with features)
   ```

2. **Small Team (3 devs)**
   ```
   Monthly Usage:
   - 300 commits (~450K tokens)
   - 60 features (~120K tokens)
   - 150 tasks (~150K tokens)
   
   OpenAI Direct Cost: ~$43.20
   YOLO Developer Package: $99 (significant savings)
   ```

3. **Development Team (10 devs)**
   ```
   Monthly Usage:
   - 1000 commits (~1.5M tokens)
   - 200 features (~400K tokens)
   - 500 tasks (~500K tokens)
   
   OpenAI Direct Cost: ~$144
   YOLO Team Package: $299 (best value)
   ```

## Enterprise Solutions

### Custom Packages
- Volume-based pricing
- Custom feature development
- On-premises deployment
- Air-gapped solutions

### Enterprise Features
1. **Security**
   - Private cloud deployment
   - Custom data handling
   - Audit logs
   - Role-based access

2. **Support**
   - 24/7 dedicated support
   - Training sessions
   - Custom documentation
   - Implementation assistance

3. **Integration**
   - CI/CD pipeline integration
   - Custom tool integration
   - SSO implementation
   - Custom API endpoints

### Getting Started with Enterprise
1. Contact sales@baudevs.com
2. Schedule solution architecture review
3. Receive custom quote
4. Implementation planning
5. Deployment and training

## Subscription Management

### How to Subscribe
1. Visit [yolo.baudevs.com/pricing](https://yolo.baudevs.com/pricing)
2. Choose your package
3. Create account or sign in
4. Enter payment information
5. Receive license key

### Managing Your Subscription
- Upgrade/downgrade anytime
- Cancel with 30-day notice
- Change billing information
- Add/remove team members

### Payment Methods
- Credit/Debit Cards
- PayPal
- Wire Transfer (Enterprise)
- Purchase Orders (Enterprise)

Remember: Choose the option that best fits your development workflow and budget. Start with your own API keys if you're just exploring, then upgrade to a YOLO package when you need the additional features and cost savings of bulk pricing.
