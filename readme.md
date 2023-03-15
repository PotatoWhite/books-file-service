# 준비

```shell
go mod init github.com/potatowhite/books/file-service

go get gorm.io/gorm
go get gorm.io/driver/potatowhite
go get github.com/spf13/viper
go get github.com/99designs/gqlgen
```

# autobind

```sql

CREATE TABLE public.folders (
                                id bigserial NOT NULL,
                                created_at timestamptz NULL,
                                updated_at timestamptz NULL,
                                deleted_at timestamptz NULL,
                                name text NOT NULL,
                                parent_id int8 NULL,
                                CONSTRAINT folders_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_folders_deleted_at ON public.folders USING btree (deleted_at);
CREATE INDEX idx_folders_parent_id ON public.folders USING btree (parent_id);


CREATE TABLE public.files (
                              id bigserial NOT NULL,
                              created_at timestamptz NULL,
                              updated_at timestamptz NULL,
                              deleted_at timestamptz NULL,
                              "name" text NOT NULL,
                              folder_id int8 NOT NULL,
                              "type" text NULL,
                              "extension" text NULL,
                              "size" int8 NULL,
                              modified text NULL,
                              CONSTRAINT files_pkey PRIMARY KEY (id),
                              CONSTRAINT fk_folders_files FOREIGN KEY (folder_id) REFERENCES public.folders(id)
);
CREATE INDEX idx_files_deleted_at ON public.files USING btree (deleted_at);
CREATE INDEX idx_files_folder_id ON public.files USING btree (folder_id);

```