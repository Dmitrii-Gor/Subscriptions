CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE subscriptions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               service_name TEXT NOT NULL,
                               price INT NOT NULL,
                               user_id UUID NOT NULL,
                               start_date DATE NOT NULL,
                               end_date DATE,
                               CONSTRAINT uniq_subscription UNIQUE (user_id, service_name, start_date, price)
);

CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);