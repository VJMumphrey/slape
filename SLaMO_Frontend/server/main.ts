import { Application } from "jsr:@oak/oak/application";
import { Router } from "jsr:@oak/oak/router";
import { oakCors } from "@tajpouria/cors";
import routeStaticFilesFrom from "./util/routeStaticFilesFrom.ts";
import data from "./api/data.json" with { type: "json"};

export const app = new Application();
const router = new Router();


router.get("/api/dinosaurs", (context) => {
  context.response.body = data;
});

router.get("/api/dinosaurs/:dinosaur", (context) => {
  if (!context?.params?.dinosaur) {
    context.response.body = "No dinosaur name provided.";
  }
  
  const dinosaur = data.find((item) => 
    item.name.toLowerCase() === context.params.dinosaur.toLowerCase()
);

context.response.body = dinosaur ?? "No dinosaur found.";
});

app.use(oakCors());
app.use(router.routes());
app.use(router.allowedMethods());
app.use(routeStaticFilesFrom([
  `${Deno.cwd()}/client/dist`,
  `${Deno.cwd()}/client/public`,
]));


if (import.meta.main) {
  console.log("Server listening on port http://localhost:8000");
  await app.listen({ port: 7739 });
}
