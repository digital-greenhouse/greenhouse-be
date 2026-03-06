-- ============================================================
--  green_house_db
--  Motor: MariaDB / InnoDB
--  Charset: utf8mb4
-- ============================================================

CREATE DATABASE IF NOT EXISTS green_house_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE green_house_db;

-- ============================================================
-- 1. USUARIOS
--    Roles: SUPERADMIN | OWNER | CLIENT
-- ============================================================
CREATE TABLE users (
    id            INT AUTO_INCREMENT PRIMARY KEY,
    name          VARCHAR(100)  NOT NULL,
    email         VARCHAR(150)  NOT NULL UNIQUE,
    password_hash VARCHAR(255)  NOT NULL,
    role          ENUM('SUPERADMIN', 'OWNER', 'CLIENT') DEFAULT 'CLIENT',
    phone         VARCHAR(20),
    is_active     BOOLEAN       DEFAULT TRUE,
    created_at    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),
    updated_at    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 2. PROPIEDADES / FINCAS
--    Pertenecen a un OWNER. Un OWNER puede tener N propiedades.
-- ============================================================
CREATE TABLE properties (
    id                    INT AUTO_INCREMENT PRIMARY KEY,
    owner_id              INT           NOT NULL,
    name                  VARCHAR(150)  NOT NULL,
    description           TEXT,
    address               VARCHAR(255),
    base_price_per_night  DECIMAL(10,2) NOT NULL,
    max_capacity          INT           NOT NULL,
    status                ENUM('ACTIVE', 'INACTIVE', 'MAINTENANCE') DEFAULT 'ACTIVE',
    created_at            TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),
    updated_at            TIMESTAMP     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_properties_owner  (owner_id),
    INDEX idx_properties_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 3. IMÁGENES DE PROPIEDADES
--    Almacenadas como Base64 en LONGTEXT.
--    Una propiedad puede tener N imágenes; una marcada como portada.
-- ============================================================
CREATE TABLE property_images (
    id            INT AUTO_INCREMENT PRIMARY KEY,
    property_id   INT           NOT NULL,
    image_data    LONGTEXT      NOT NULL,          -- Base64 encoded image
    mime_type     VARCHAR(50)   NOT NULL,          -- Ej: image/jpeg, image/png, image/webp
    alt_text      VARCHAR(150),                    -- Texto alternativo / descripción
    is_cover      BOOLEAN       DEFAULT FALSE,     -- Solo una debe ser TRUE por propiedad
    sort_order    INT           DEFAULT 0,         -- Para ordenar el carrusel
    created_at    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),

    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    INDEX idx_images_property (property_id),
    INDEX idx_images_cover    (property_id, is_cover)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 4. REGLAS DE PRECIOS DINÁMICOS
--    Modificadores por temporada/fecha. Ej: Semana Santa = ×1.50
--    El motor de cotización las aplica automáticamente.
-- ============================================================
CREATE TABLE pricing_rules (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    property_id     INT           NOT NULL,
    name            VARCHAR(100)  NOT NULL,         -- Ej: "Temporada Alta Diciembre"
    start_date      DATE          NOT NULL,
    end_date        DATE          NOT NULL,
    price_modifier  DECIMAL(5,2)  NOT NULL,         -- 1.50 = +50% | 0.80 = -20%
    description     VARCHAR(255),
    is_active       BOOLEAN       DEFAULT TRUE,
    created_at      TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),

    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    INDEX idx_pricing_property (property_id),
    INDEX idx_pricing_dates    (property_id, start_date, end_date),

    CONSTRAINT chk_pricing_dates CHECK (end_date >= start_date),
    CONSTRAINT chk_modifier_positive CHECK (price_modifier > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 5. COTIZACIONES (Inteligencia Proyectiva)
--    Registro de TODA intención de reserva, incluso abandonadas.
--    Permite analítica de conversión y embudo de ventas.
--    client_id puede ser NULL (usuario invitado/anónimo).
-- ============================================================
CREATE TABLE quotes (
    id                    INT AUTO_INCREMENT PRIMARY KEY,
    property_id           INT           NOT NULL,
    client_id             INT           NULL,        -- NULL = invitado anónimo
    check_in_date         DATE          NOT NULL,
    check_out_date        DATE          NOT NULL,
    guest_count           INT           NOT NULL,
    calculated_total      DECIMAL(10,2) NOT NULL,    -- Total calculado con reglas de precio
    nights_count          INT           NOT NULL,    -- (check_out - check_in) en días
    applied_modifier      DECIMAL(5,2)  DEFAULT 1.00, -- Modificador que se aplicó
    status                ENUM('ACTIVE','CONVERTED','EXPIRED','ABANDONED') DEFAULT 'ACTIVE',
    abandonment_reason    VARCHAR(255)  NULL,
    expires_at            TIMESTAMP     NULL,         -- Validez de la cotización
    created_at            TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),
    updated_at            TIMESTAMP     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id)   REFERENCES users(id)      ON DELETE SET NULL,

    INDEX idx_quotes_property (property_id),
    INDEX idx_quotes_client   (client_id),
    INDEX idx_quotes_status   (status),
    INDEX idx_quotes_dates    (property_id, check_in_date, check_out_date),

    CONSTRAINT chk_quotes_dates      CHECK (check_out_date > check_in_date),
    CONSTRAINT chk_quotes_guests     CHECK (guest_count > 0),
    CONSTRAINT chk_quotes_total      CHECK (calculated_total >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 6. RESERVAS (Núcleo operativo)
--    Puede nacer de una quote (quote_id) o crearse directamente.
--    Ciclo de vida: PENDING_PAYMENT → CONFIRMED → COMPLETED
--                                  ↘ CANCELLED
-- ============================================================
CREATE TABLE bookings (
    id                  INT AUTO_INCREMENT PRIMARY KEY,
    property_id         INT           NOT NULL,
    client_id           INT           NOT NULL,
    quote_id            INT           NULL,          -- Cotización de origen (opcional)
    check_in_date       DATE          NOT NULL,
    check_out_date      DATE          NOT NULL,
    guest_count         INT           NOT NULL,
    nights_count        INT           NOT NULL,
    total_price         DECIMAL(10,2) NOT NULL,
    status              ENUM('PENDING_PAYMENT','CONFIRMED','CANCELLED','COMPLETED')
                            DEFAULT 'PENDING_PAYMENT',
    cancellation_reason VARCHAR(255)  NULL,          -- Solo si status = CANCELLED
    special_requests    TEXT          NULL,           -- Notas del cliente
    created_at          TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),
    updated_at          TIMESTAMP     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id)   REFERENCES users(id)      ON DELETE CASCADE,
    FOREIGN KEY (quote_id)    REFERENCES quotes(id)     ON DELETE SET NULL,

    -- Índice crítico para consultas de disponibilidad
    INDEX idx_bookings_availability (property_id, check_in_date, check_out_date),
    INDEX idx_bookings_client       (client_id),
    INDEX idx_bookings_status       (status),

    CONSTRAINT chk_bookings_dates  CHECK (check_out_date > check_in_date),
    CONSTRAINT chk_bookings_guests CHECK (guest_count > 0),
    CONSTRAINT chk_bookings_price  CHECK (total_price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 7. PAGOS / COMPROBANTES
--    Una reserva puede tener múltiples pagos (abono + saldo).
--    proof_data almacena el comprobante en Base64 si aplica.
-- ============================================================
CREATE TABLE payments (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    booking_id      INT           NOT NULL,
    amount          DECIMAL(10,2) NOT NULL,
    payment_method  ENUM('TRANSFERENCIA','EFECTIVO','PASARELA') NOT NULL,
    proof_data      LONGTEXT      NULL,              -- Base64 del comprobante (imagen/PDF)
    proof_mime_type VARCHAR(50)   NULL,              -- Ej: image/jpeg, application/pdf
    status          ENUM('PENDING_VERIFICATION','VERIFIED','REJECTED')
                        DEFAULT 'PENDING_VERIFICATION',
    rejection_reason VARCHAR(255) NULL,              -- Solo si status = REJECTED
    verified_by     INT           NULL,              -- ID del admin/owner que verificó
    payment_date    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),
    verified_at     TIMESTAMP     NULL,

    FOREIGN KEY (booking_id)   REFERENCES bookings(id) ON DELETE CASCADE,
    FOREIGN KEY (verified_by)  REFERENCES users(id)    ON DELETE SET NULL,

    INDEX idx_payments_booking (booking_id),
    INDEX idx_payments_status  (status),

    CONSTRAINT chk_payment_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- 8. RESEÑAS / CALIFICACIONES (Opcional pero recomendado)
--    Solo clientes con reserva COMPLETED pueden dejar reseña.
-- ============================================================
CREATE TABLE reviews (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    property_id INT           NOT NULL,
    client_id   INT           NOT NULL,
    booking_id  INT           NOT NULL UNIQUE,  -- 1 reseña por reserva completada
    rating      TINYINT       NOT NULL,         -- 1 a 5 estrellas
    comment     TEXT          NULL,
    is_visible  BOOLEAN       DEFAULT TRUE,
    created_at  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP(),

    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id)   REFERENCES users(id)      ON DELETE CASCADE,
    FOREIGN KEY (booking_id)  REFERENCES bookings(id)   ON DELETE CASCADE,

    INDEX idx_reviews_property (property_id),

    CONSTRAINT chk_rating_range CHECK (rating BETWEEN 1 AND 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================
-- FIN DEL ESQUEMA
-- ============================================================