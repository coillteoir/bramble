import express, { Express, Request, Response } from "express";
import cors from "cors";
import { getPL, getPo } from "./query_cluster";

const app: Express = express();

const port: number = 5555;

app.use(cors());
app.use(express.json());

app.get("/", async (req: Request, res: Response) => {
  try {
    const Pipelines = await getPL("default");
    console.log(Pipelines);
    res.json(Pipelines);
  } catch (error) {
    res.status(500).json({ error: "internal server error (womp womp)" });
  }
});

app.listen(port, () => {
  console.log("Server is running");
});
