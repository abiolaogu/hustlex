import type { AuthBindings } from "@refinedev/core";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8081";

export const authProvider: AuthBindings = {
  login: async ({ email, password }) => {
    try {
      const response = await fetch(`${API_URL}/api/v1/auth/admin/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem("auth", JSON.stringify(data));
        return {
          success: true,
          redirectTo: "/",
        };
      }

      return {
        success: false,
        error: {
          name: "Login Error",
          message: "Invalid email or password",
        },
      };
    } catch (error) {
      return {
        success: false,
        error: {
          name: "Login Error",
          message: "An error occurred during login",
        },
      };
    }
  },

  logout: async () => {
    localStorage.removeItem("auth");
    return {
      success: true,
      redirectTo: "/login",
    };
  },

  check: async () => {
    const auth = localStorage.getItem("auth");
    if (auth) {
      return {
        authenticated: true,
      };
    }

    return {
      authenticated: false,
      redirectTo: "/login",
    };
  },

  getPermissions: async () => {
    const auth = localStorage.getItem("auth");
    if (auth) {
      const { role } = JSON.parse(auth);
      return role;
    }
    return null;
  },

  getIdentity: async () => {
    const auth = localStorage.getItem("auth");
    if (auth) {
      const { user } = JSON.parse(auth);
      return user;
    }
    return null;
  },

  onError: async (error) => {
    if (error?.status === 401) {
      return {
        logout: true,
        redirectTo: "/login",
      };
    }
    return { error };
  },
};
