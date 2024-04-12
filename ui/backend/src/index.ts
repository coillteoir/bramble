import express, { Express, Request, Response } from "express";
import cors from "cors";
import { getPipelines, getExecutions, getJobs } from "./query_cluster";
import { getLogger } from "log4js";
const app: Express = express();

const port: number = process.env.PORT ? parseInt(process.env.PORT) : NaN;
const logger = getLogger();
logger.level = "info";

app.use(cors());
app.use(express.json());
app.use(express.static("public"));

app.get(
  "/resources/:resource/:namespace",
  async (req: Request, res: Response) => {
    logger.info(`querying ${req.params.resource} in ${req.params.namespace}`);
    try {
      switch (req.params.resource) {
        case "jobs": {
          const jobs = await getJobs(req.params.namespace);
          res.json(jobs);
          break;
        }
        case "pipelines": {
          const pipelines = await getPipelines(req.params.namespace);
          res.json(pipelines);
          break;
        }
        case "executions": {
          const executions = await getExecutions(req.params.namespace);
          res.json(executions);
          break;
        }
        default:
          logger.error(`Cannot query unknown resource: ${req.params.resource}`);
      }
    } catch (error) {
      res.status(500).json({
        error: `Cannnot fetch ${req.params.resource} from namespace: ${req.params.namespace}`,
      });
    }
  },
);
if (!isNaN(port)) {
  app.listen(port, () => {
    console.log("Server is running");
  });
}
