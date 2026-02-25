import {
  createAuthServiceClient,
  createUserServiceClient,
  type AuthService,
  type UserService,
} from "../gen/krathub/service/v1";

import {
  createRequestHandler,
  type RequestHandlerOptions,
} from "./requestHandler";

export type KrathubClients = {
  auth: AuthService;
  user: UserService;
};

export function createKrathubClients(
  options: RequestHandlerOptions = {},
): KrathubClients {
  const handler = createRequestHandler(options);
  return {
    auth: createAuthServiceClient(handler),
    user: createUserServiceClient(handler),
  };
}

export * from "../gen/krathub/service/v1";
