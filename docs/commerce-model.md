# Commerce Model

Phase 5 keeps the commerce world intentionally small.

## Entities

- `Buyer`: identity and lightweight risk tier
- `Merchant`: merchant identity and category
- `Product`: one purchasable item with amount and currency
- `Order`: the central execution entity, including actor and delegation fields
- `Payment`: an authorization outcome for an order
- `Refund`: a post-purchase action used for the suspicious baseline flow

## Design rule

These entities exist to support trust-lab execution, not to become a generic ecommerce backend.

That is why the model focuses on:

- submitted actor identity
- delegation mode
- mandate and provenance references
- event emission
- trust-decision outcomes

It intentionally does not model:

- inventory systems
- shipping
- taxes
- gateway integrations
- pricing engines
