const browserWindow: any = window || {};
const browserWindowEnv = browserWindow['env'] || {};

export const environment = {
  isProduction: true,
  apiUrl: browserWindowEnv.apiUrl,
  autoUpdateIssuesMin: 2
};
