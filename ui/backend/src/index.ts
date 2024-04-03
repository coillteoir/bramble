import express, { Express, Request, Response } from "express";
import cors from "cors";
import { getPipelines, getExecutions, getPods } from "./query_cluster";

const app: Express = express();

const port: number = 5555;

app.use(cors());
app.use(express.json());
app.use(express.static("public"));

app.get(
  "/resources/:resource/:namespace",
  async (req: Request, res: Response) => {
    console.log(`querying ${req.params.resource} in ${req.params.namespace}`);
    try {
      switch (req.params.resource) {
        case "pods": {
          const pods = await getPods(req.params.namespace);
          console.log(pods);
          res.json(pods);
          break;
        }
        case "pipelines": {
          const pipelines = await getPipelines(req.params.namespace);
          console.log(pipelines);
          res.json(pipelines);
          break;
        }
        case "executions": {
          const executions = await getExecutions(req.params.namespace);
          console.log(executions);
          res.json(executions);
          break;
        }
        default:
          console.log(`Cannot query unknown resource: ${req.params.resource}`);
      }
    } catch (error) {
      res.status(500).json({
        error: `Cannnot fetch ${req.params.resource} from namespace: ${req.params.namespace}`,
      });
    }
  },
);

app.listen(port, () => {
  console.log("Server is running");
});
