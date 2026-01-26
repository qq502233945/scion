/**
 * Authentication middleware
 *
 * Provides route protection and user context for authenticated routes.
 * Works alongside dev-auth middleware - dev-auth sets user if dev token is present,
 * this middleware enforces authentication for protected routes.
 */

import type { Context, Next } from 'koa';

import type { User } from '../../shared/types.js';
import type { AppConfig } from '../config.js';

/**
 * Extended Koa state with auth information
 */
export interface AuthState {
  /** Currently authenticated user */
  user?: User;
  /** Dev token (if dev auth is enabled) */
  devToken?: string;
  /** Whether dev auth is enabled */
  devAuthEnabled?: boolean;
  /** Request ID for tracing */
  requestId?: string;
}

/**
 * Check if a URL should be protected (require authentication)
 */
function isProtectedRoute(url: string): boolean {
  // Public routes that don't require authentication
  const publicPaths = [
    '/healthz',
    '/readyz',
    '/login',
    '/auth/login',
    '/auth/callback',
    '/auth/error',
    '/assets/',
    '/favicon.ico',
  ];

  // Check if the URL matches any public path
  for (const path of publicPaths) {
    if (url.startsWith(path)) {
      return false;
    }
  }

  return true;
}

/**
 * Check if the request accepts HTML (is a browser request)
 */
function acceptsHtml(ctx: Context): boolean {
  const accept = ctx.headers.accept || '';
  return accept.includes('text/html');
}

/**
 * Create auth middleware that enforces authentication on protected routes
 *
 * @param _config - Application configuration (reserved for future use)
 * @returns Koa middleware function
 */
export function createAuthMiddleware(_config: AppConfig) {
  return async function authMiddleware(ctx: Context, next: Next) {
    const state = ctx.state as AuthState;

    // Check if route needs protection
    if (!isProtectedRoute(ctx.path)) {
      return next();
    }

    // Check if user is authenticated (either from dev-auth or session)
    let user: User | undefined = state.user;

    // If no user from dev-auth, check session
    if (!user && ctx.session?.user) {
      user = ctx.session.user;
      if (user) {
        state.user = user;
      }
    }

    // If still no user, handle unauthenticated request
    if (!user) {
      // Store the original URL for redirect after login
      if (ctx.session) {
        ctx.session.returnTo = ctx.originalUrl;
      }

      // For API requests, return 401
      if (ctx.path.startsWith('/api/')) {
        ctx.status = 401;
        ctx.body = {
          error: 'Unauthorized',
          message: 'Authentication required',
        };
        return;
      }

      // For browser requests, redirect to login
      if (acceptsHtml(ctx)) {
        ctx.redirect('/auth/login');
        return;
      }

      // Default: 401 response
      ctx.status = 401;
      ctx.body = {
        error: 'Unauthorized',
        message: 'Authentication required',
      };
      return;
    }

    // User is authenticated - continue
    await next();
  };
}

/**
 * Validate that a user's email domain is authorized
 *
 * @param email - User email address
 * @param authorizedDomains - List of authorized email domains
 * @returns true if authorized, false otherwise
 */
export function isEmailAuthorized(email: string, authorizedDomains: string[]): boolean {
  // If no domains are configured, allow all
  if (!authorizedDomains || authorizedDomains.length === 0) {
    return true;
  }

  // Extract domain from email
  const atIndex = email.lastIndexOf('@');
  if (atIndex === -1) {
    return false;
  }

  const domain = email.substring(atIndex + 1).toLowerCase();

  // Check if domain is in the authorized list
  return authorizedDomains.some((authorized) => authorized.toLowerCase() === domain);
}
