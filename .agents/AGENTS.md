# Kirmya Platform Implementation Rules

## Monetization & Billing Restriction (Current Phase)
Kirmya is currently a completely free platform.
Do **not** implement or display any subscription, premium plans, payment gateways, billing pages, invoices, coupons, payment methods, pricing management, or monetization features.

Until explicitly enabled in the future:
* Hide all Premium, Pro, Business, Enterprise, Upgrade, and Subscription UI.
* Do not create Billing or Subscription pages in the user settings.
* Do not add payment APIs, database tables, or payment providers (Stripe, PayPal, Razorpay, etc.).
* Do not show locked features or upsell banners.
* Every feature available in the current platform should be accessible to all authenticated users without payment.

Design the architecture to be modular and future-ready so that subscriptions can be integrated later without major refactoring, but keep all payment-related functionality disabled and excluded from the current implementation.

Replace the "Billing & Subscription" section in Settings with a simple **Platform Information** section containing:
* Current Platform Version
* Release Notes
* What's New
* Feature Roadmap (optional)
* Open Source Libraries Used (optional)
* Legal Information
* Terms of Service
* Privacy Policy
* Cookie Policy
* Licenses & Attributions
