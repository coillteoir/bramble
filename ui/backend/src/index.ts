import express, { Express, Request, Response } from "express";
import cors from "cors";
import { getPipelines } from "./query_cluster";

const app: Express = express();

const port: number = 5555;

app.use(cors());
app.use(express.json());
app.use(express.static("public"));

app.get("/pipelines/:ns", async (req: Request, res: Response) => {
  try {
    console.log("Queriying pipelines in: " + req.params.ns);
    const Pipelines = await getPipelines(req.params.ns);
    console.log(Pipelines);
    res.json(Pipelines);
  } catch (error) {
    console.log(error);
    res.status(500).json({ error: "internal server error (womp womp)" });
  }
});

app.listen(port, () => {
  console.log("Server is running");
});
