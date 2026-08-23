export default function sitemap() {
  const baseUrl = "https://gitresolve.dev";

  const routes = ["", "/get-started", "/architecture", "/merge-flow", "/commands", "/operations"];

  return [
    ...routes.map((route, index) => ({
      url: `${baseUrl}${route}`,
      lastModified: new Date(),
      changeFrequency: "weekly",
      priority: index === 0 ? 1 : 0.8,
    })),
  ];
}
