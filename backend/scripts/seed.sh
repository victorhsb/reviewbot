#!/bin/bash

psql "postgres://user:password@localhost:5432/product-reviews" --file ./scripts/seed.sql
