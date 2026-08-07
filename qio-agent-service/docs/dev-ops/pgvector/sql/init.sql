CREATE EXTENSION IF NOT EXISTS vector;

-- 查询表；SELECT * FROM information_schema.tables

-- 删除旧的表（如果存在，CASCADE 确保删除所有依赖）
DROP TABLE IF EXISTS public.vector_store CASCADE;

-- 创建新的表，使用UUID作为主键
CREATE TABLE public.vector_store (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    metadata JSONB,
    embedding VECTOR(1536)
);

-- 删除旧的表（如果存在，CASCADE 确保删除所有依赖）
DROP TABLE IF EXISTS public.store_openai CASCADE;

-- 创建新的表，使用UUID作为主键
CREATE TABLE public.store_openai (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    metadata JSONB,
    embedding VECTOR(1536)
);

-- 删除旧的表（如果存在，CASCADE 确保删除所有依赖）
DROP TABLE IF EXISTS public.vector_store_openai CASCADE;

-- 创建新的表，使用UUID作为主键，embedding 维度为 1536（适配 text-embedding-v4 模型实际输出）
CREATE TABLE public.vector_store_openai (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    metadata JSONB,
    embedding VECTOR(1536)
);
