import {
  createAuthServiceClient,
  createUserServiceClient,
  type AuthService,
  type UserService,
} from "../gen/micro_forge/service/v1";

import {
  createRequestHandler,
  type RequestHandlerOptions,
} from "./requestHandler";

export type microforgeClients = {
  auth: AuthService;
  user: UserService;
};

export function createmicroforgeClients(
  options: RequestHandlerOptions = {},
): microforgeClients {
  const handler = createRequestHandler(options);
  return {
    auth: createAuthServiceClient(handler),
    user: createUserServiceClient(handler),
  };
}

export * from "../gen/micro_forge/service/v1";
