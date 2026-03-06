const browserWindow: any = window || {};
const browserWindowEnv = browserWindow['env'] || {};

export const environment = {
  isProduction: true,
  apiUrl: browserWindowEnv.apiUrl || `${window.location.origin}/api`,
  autoUpdateIssuesMin: 2
};
