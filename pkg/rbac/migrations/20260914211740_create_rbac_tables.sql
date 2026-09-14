-- +goose Up

CREATE TABLE public.actions (
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT actions_pkey PRIMARY KEY (name)
);

CREATE TABLE public.roles (
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT roles_pkey PRIMARY KEY (name)
);

CREATE TABLE public.resources (
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT resources_pkey PRIMARY KEY (name)
);

CREATE TABLE public.permissions (
    resource_name TEXT NOT NULL,
    action_name TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT permissions_pkey PRIMARY KEY (resource_name, action_name),

    CONSTRAINT fk_permissions_action
        FOREIGN KEY (action_name)
        REFERENCES public.actions(name)
        ON DELETE CASCADE,

    CONSTRAINT fk_permissions_resource
        FOREIGN KEY (resource_name)
        REFERENCES public.resources(name)
        ON DELETE CASCADE
);

CREATE TABLE public.role_permissions (
    role_name TEXT NOT NULL,
    resource_name TEXT NOT NULL,
    action_name TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT role_permissions_pkey
        PRIMARY KEY (role_name, resource_name, action_name),

    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (resource_name, action_name)
        REFERENCES public.permissions(resource_name, action_name)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY (role_name)
        REFERENCES public.roles(name)
        ON DELETE CASCADE
);

CREATE TABLE public.user_roles (
    user_id TEXT NOT NULL,
    role_name TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT user_roles_pkey PRIMARY KEY (user_id, role_name),

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_name)
        REFERENCES public.roles(name)
        ON DELETE CASCADE
);

-- Initial data

INSERT INTO public.actions (name, created_at, updated_at)
VALUES
    ('write', NOW(), NOW()),
    ('read', NOW(), NOW());

INSERT INTO public.roles (name, description, created_at, updated_at)
VALUES
    ('user', 'user', NOW(), NOW()),
    ('admin', 'admin', NOW(), NOW()),
    ('superadmin', 'superadmin', NOW(), NOW());

INSERT INTO public.user_roles (
    user_id,
    role_name,
    created_at,
    updated_at
)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'superadmin',
    NOW(),
    NOW()
);


-- +goose Down

DROP TABLE IF EXISTS public.user_roles;
DROP TABLE IF EXISTS public.role_permissions;
DROP TABLE IF EXISTS public.permissions;
DROP TABLE IF EXISTS public.resources;
DROP TABLE IF EXISTS public.roles;
DROP TABLE IF EXISTS public.actions;