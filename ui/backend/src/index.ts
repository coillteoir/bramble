import express, { Express, Request, Response } from 'express'
import cors from 'cors'

const k8s = require('@kubernetes/client-node');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();

const k8sApi = kc.makeApiClient(k8s.CoreV1Api);

const getPo = async (ns: string) => {
       const response =  await k8sApi.listNamespacedPod(ns)
       const podNames = response.body.items.map(pod => pod.metadata.name)
       console.log(podNames)
       return podNames
};

const app: Express = express()

const port: number = 5555

app.use(cors())
app.use(express.json())

app.get('/', async (req: Request, res: Response) => {
    try {
        const podNames = await getPo('default')
        res.json(podNames)
    } catch (error) {
        res.status(500).json({error: "internal server error (womp womp)"})
    }
})

app.listen(port, () => {
    console.log("Server is running")
})
