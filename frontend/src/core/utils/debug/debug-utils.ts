import { logger } from '../logger/logger';

/**
 * Debug utility for API requests
 */
export const apiDebug = {
  logRequest: (config: any) => {
    logger.debug('🚀 API Request:', {
      url: config.url,
      method: config.method,
      headers: config.headers,
      data: config.data
    });
  },

  logResponse: (response: any) => {
    logger.debug('✅ API Response:', {
      status: response.status,
      data: response.data,
      headers: response.headers
    });
  },

  logError: (error: any) => {
    logger.error('❌ API Error:', {
      message: error.message,
      status: error.response?.status,
      data: error.response?.data,
      config: error.config
    });
  }
};

/**
 * Debug utility for React components
 */
export const componentDebug = {
  logRender: (componentName: string, props: any) => {
    logger.debug(`🔄 ${componentName} Rendering:`, {
      props,
      timestamp: new Date().toISOString()
    });
  },

  logState: (componentName: string, state: any) => {
    logger.debug(`📊 ${componentName} State:`, {
      state,
      timestamp: new Date().toISOString()
    });
  }
};

/**
 * Debug utility for hooks
 */
export const hookDebug = {
  logHookCall: (hookName: string, params: any) => {
    logger.debug(`🎣 ${hookName} Called:`, {
      params,
      timestamp: new Date().toISOString()
    });
  },

  logHookResult: (hookName: string, result: any) => {
    logger.debug(`🎣 ${hookName} Result:`, {
      result,
      timestamp: new Date().toISOString()
    });
  }
};