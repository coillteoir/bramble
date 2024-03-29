import express, { Express, Request, Response } from "express";
import cors from "cors";
import { getPipelines, getExecutions, getPods } from "./query_cluster";

const app: Express = express();

const port: number = 5555;

app.use(cors());
app.use(express.json());
app.use(express.static("public"));

app.get("/pipelines/:ns", async (req: Request, res: Response) => {
  try {
    console.log("Querying pipelines in:" + req.params.ns);
    const pipelines = await getPipelines(req.params.ns);
    console.log(pipelines);
    res.json(pipelines);
  } catch (error) {
    console.log(error);
    res.status(500).json({
      error: "Cannot fetch pipelines from namespace: " + req.params.ns,
    });
  }
});

app.get("/pods/:ns", async (req: Request, res: Response) => {
  try {
    console.log("Querying pods in:" + req.params.ns);
    const pods = await getPods(req.params.ns);
    console.log(pods);
    res.json(pods);
  } catch (error) {
    console.log(error);
    res
      .status(500)
      .json({ error: "Cannot fetch pods from namespace: " + req.params.ns });
  }
});


app.get("/executions/:ns", async (req: Request, res: Response) => {
try {
    console.log("Querying executions in:" + req.params.ns);
    const executions = await getExecutions(req.params.ns);
    console.log(executions);
    res.json(executions);
  } catch (error) {
    console.log(error);
    res
      .status(500)
      .json({ error: "Cannot fetch executions from namespace: " + req.params.ns });
  }
})

app.listen(port, () => {
  console.log("Server is running");
});
