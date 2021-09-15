export const routeUrls = {
  login: "/login",
  pictures: "/",
  picture: {
    route: "/:id",
    link: (id: string) => `/${id}`,
  },
  analytics: "/analytics",
};
