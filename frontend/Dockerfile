FROM node:23-alpine as base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY ./frontend /app
COPY ./frontend/../.env.docker /.env.public
WORKDIR /app

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM base
COPY --from=build /app/node_modules /app/node_modules
COPY --from=build /app/dist /app/dist

EXPOSE 80

CMD [ "pnpm", "run", "preview" ]
