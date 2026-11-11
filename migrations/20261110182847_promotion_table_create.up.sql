CREATE TYPE PROMOTIONS_DISCOUNT_TYPE_ENUM AS ENUM (
    'PERCENTAGE',
    'FIXED'
    );

CREATE TYPE PROMOTIONS_STATUS_ENUM AS ENUM (
    'ACTIVE',
    'DISABLED'
    );

CREATE TYPE PROMO_REDEMPTIONS_STATUS_ENUM AS ENUM (
    'RESERVED',
    'REDEEMED',
    'RELEASED'
    );

CREATE TABLE IF NOT EXISTS promotions
(
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    discount_type PROMOTIONS_DISCOUNT_TYPE_ENUM NOT NULL,
    discount_value INT NOT NULL,
    max_discount INT NOT NULL,
    max_uses BIGINT NOT NULL,
    used_count BIGINT NOT NULL DEFAULT 0,
    starts_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    status PROMOTIONS_STATUS_ENUM NOT NULL
);

CREATE TABLE IF NOT EXISTS promo_codes
(
    id UUID PRIMARY KEY,
    code VARCHAR NOT NULL,
    promotion_id UUID NOT NULL,
    max_uses_per_user BIGINT NOT NULL,

    CONSTRAINT fk_promo_codes_promotion
        FOREIGN KEY (promotion_id)
            REFERENCES promotions(id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS promo_redemptions
(
    id UUID PRIMARY KEY,
    promotion_id UUID NOT NULL,
    user_id UUID NOT NULL,
    order_id UUID NOT NULL,
    discount_amount BIGINT NOT NULL,
    status PROMO_REDEMPTIONS_STATUS_ENUM NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_promo_redemptions_promotion
        FOREIGN KEY (promotion_id)
            REFERENCES promotions(id)
            ON DELETE RESTRICT,

    CONSTRAINT fk_promo_redemptions_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE RESTRICT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_promo_code
    ON promo_codes(code);

CREATE UNIQUE INDEX IF NOT EXISTS ids_promo_order
    ON promo_redemptions(order_id, promotion_id);

CREATE INDEX IF NOT EXISTS idx_user_redemptions
    ON promo_redemptions(user_id, promotion_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_promo
    ON promo_redemptions(user_id, promotion_id)
    WHERE status = 'REDEEMED';
