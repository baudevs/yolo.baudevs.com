# YOLO CLI Webhook Handler

This directory contains the webhook handler for Stripe events, designed to be deployed on Vercel.

## Setup

1. Deploy to Vercel:
```bash
vercel
```

2. Add environment variables in Vercel:
```bash
vercel env add STRIPE_SECRET_KEY
vercel env add STRIPE_WEBHOOK_SECRET
```

3. Configure webhook in Stripe:
- Go to Stripe Dashboard > Developers > Webhooks
- Add endpoint URL: https://your-vercel-domain.vercel.app/api/stripe
- Select events to listen for:
  - customer.subscription.created
  - customer.subscription.updated
  - customer.subscription.deleted
  - invoice.paid
  - invoice.payment_failed

4. Get the webhook signing secret from Stripe and add it to Vercel:
```bash
vercel env add STRIPE_WEBHOOK_SECRET
```

## Development

To test locally:
1. Install Stripe CLI
2. Forward webhooks:
```bash
stripe listen --forward-to localhost:3000/api/stripe
```

## Event Handling

The webhook handles these Stripe events:

- `customer.subscription.created`: When a new subscription is created
- `customer.subscription.updated`: When a subscription status changes
- `customer.subscription.deleted`: When a subscription is cancelled
- `invoice.paid`: When a payment succeeds
- `invoice.payment_failed`: When a payment fails

## Security

- All requests are verified using Stripe's webhook signature
- Environment variables are securely stored in Vercel
- Only POST requests are allowed
- Response headers are properly set

## Logging

Events are logged using Vercel's built-in logging system. View logs in the Vercel dashboard under:
- Project > Deployments > Select deployment > View Logs
