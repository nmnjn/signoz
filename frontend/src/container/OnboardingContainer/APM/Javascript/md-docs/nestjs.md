## Requirements

**Supported Versions**

- `>=4.0.0`

## Send traces to SigNoz Cloud

Based on your application environment, you can choose the setup below to send traces to SigNoz Cloud.

### Application on VMs

From VMs, there are two ways to send data to SigNoz Cloud.

- Send traces directly to SigNoz Cloud (quick start)
- Send traces via OTel Collector binary (recommended)

#### **Send traces directly to SigNoz Cloud**

Step 1. Install OpenTelemetry packages

```bash
npm install --save @opentelemetry/api@^1.4.1                                                                       
npm install --save @opentelemetry/sdk-node@^0.39.1
npm install --save @opentelemetry/auto-instrumentations-node@^0.37.0
npm install --save @opentelemetry/exporter-trace-otlp-http@^0.39.1
```

Step 2. Create `tracer.ts` file

You need to configure the endpoint for SigNoz cloud in this file. You also need to configure your service name. In this example, we have used `sampleNestjsApplication`.

```js
'use strict'
const process = require('process');
const opentelemetry = require('@opentelemetry/sdk-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');
const {Resource} = require('@opentelemetry/resources');
const {SemanticResourceAttributes} = require('@opentelemetry/semantic-conventions');

const exporterOptions = {
    url: 'https://ingest.{region}.signoz.cloud:443/v1/traces'
  }

const traceExporter = new OTLPTraceExporter(exporterOptions);
const sdk = new opentelemetry.NodeSDK({
  traceExporter,
  instrumentations: [getNodeAutoInstrumentations()],
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'sampleNestjsApplication'
  })
  });
  
  // initialize the SDK and register with the OpenTelemetry API
  // this enables the API to record telemetry
  sdk.start()
  
  // gracefully shut down the SDK on process exit
  process.on('SIGTERM', () => {
    sdk.shutdown()
    .then(() => console.log('Tracing terminated'))
    .catch((error) => console.log('Error terminating tracing', error))
    .finally(() => process.exit(0));
    });
    
  module.exports = sdk
```

Depending on the choice of your region for SigNoz cloud, the ingest endpoint will vary accordingly.

 US -	ingest.us.signoz.cloud:443/v1/traces 

 IN -	ingest.in.signoz.cloud:443/v1/traces 

 EU - ingest.eu.signoz.cloud:443/v1/traces 



Step 3. Import the tracer module where your app starts  `(Ex —> main.ts)`
    
```jsx
const tracer = require('./tracer')
```
    

Step 4. Start the tracer

In the `async function boostrap` section of the application code `(Ex —> In main.ts)`, initialize the tracer as follows: 

```jsx
const tracer = require('./tracer')

import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
  // All of your application code and any imports that should leverage
  // OpenTelemetry automatic instrumentation must go here.

async function bootstrap() {
    await tracer.start();
    const app = await NestFactory.create(AppModule);
    await app.listen(3001);
  }
  bootstrap();
```

Step 5. Run the application

```bash
OTEL_EXPORTER_OTLP_HEADERS="signoz-access-token=<SIGNOZ_INGESTION_KEY>" nest start
```

You can now run your Nestjs application. The data captured with OpenTelemetry from your application should start showing on the SigNoz dashboard.

`<SIGNOZ_INGESTION_KEY>` is the API token provided by SigNoz. You can find your ingestion key from SigNoz cloud account details sent on your email.

---

#### **Send traces via OTel Collector binary**

OTel Collector binary helps to collect logs, hostmetrics, resource and infra attributes. It is recommended to install Otel Collector binary to collect and send traces to SigNoz cloud. You can correlate signals and have rich contextual data through this way.

You can find instructions to install OTel Collector binary [here](https://signoz.io/docs/tutorial/opentelemetry-binary-usage-in-virtual-machine/) in your VM. Once you are done setting up your OTel Collector binary, you can follow the below steps for instrumenting your Javascript application.

Step 1. Install OpenTelemetry packages

```js
npm install --save @opentelemetry/api@^1.4.1
npm install --save @opentelemetry/sdk-node@^0.39.1
npm install --save @opentelemetry/auto-instrumentations-node@^0.37.0
npm install --save @opentelemetry/exporter-trace-otlp-http@^0.39.1
```

Step 2. Create `tracer.ts` file

You need to configure your service name. In this example, we have used `sampleNestjsApplication`.

```js
'use strict'
const process = require('process');
//OpenTelemetry
const opentelemetry = require('@opentelemetry/sdk-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');
const {Resource} = require('@opentelemetry/resources');
const {SemanticResourceAttributes} = require('@opentelemetry/semantic-conventions');

const exporterOptions = {
    url: 'http://localhost:4318/v1/traces'
  }

const traceExporter = new OTLPTraceExporter(exporterOptions);
const sdk = new opentelemetry.NodeSDK({
  traceExporter,
  instrumentations: [getNodeAutoInstrumentations()],
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'sampleNestjsApplication'
  })
  });
  
  // initialize the SDK and register with the OpenTelemetry API
  // this enables the API to record telemetry
  sdk.start()
  
  // gracefully shut down the SDK on process exit
  process.on('SIGTERM', () => {
    sdk.shutdown()
    .then(() => console.log('Tracing terminated'))
    .catch((error) => console.log('Error terminating tracing', error))
    .finally(() => process.exit(0));
    });
    
  module.exports = sdk
```

Step 3. Import the tracer module where your app starts
    
```jsx
const tracer = require('./tracer')
```
    

Step 4. Start the tracer

In the `async function boostrap` section of the application code, initialize the tracer as follows: 

```jsx
const tracer = require('./tracer')

import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
  // All of your application code and any imports that should leverage
  // OpenTelemetry automatic instrumentation must go here.

async function bootstrap() {
    await tracer.start();
    const app = await NestFactory.create(AppModule);
    await app.listen(3001);
  }
  bootstrap();
```

Step 5. Run the application

---

### Applications Deployed on Kubernetes

For Javascript application deployed on Kubernetes, you need to install OTel Collector agent in your k8s infra to collect and send traces to SigNoz Cloud. You can find the instructions to install OTel Collector agent [here](https://signoz.io/docs/tutorial/kubernetes-infra-metrics/).

Once you have set up OTel Collector agent, you can proceed with OpenTelemetry Javascript instrumentation by following the below steps:

Step 1. Install OpenTelemetry packages

```bash
npm install --save @opentelemetry/api@^1.4.1
npm install --save @opentelemetry/sdk-node@^0.39.1
npm install --save @opentelemetry/auto-instrumentations-node@^0.37.0
npm install --save @opentelemetry/exporter-trace-otlp-http@^0.39.1
```

Step 2. Create `tracer.ts` file

You need to configure your service name. In this example, we have used `sampleNestjsApplication`.

```js
'use strict'
const process = require('process');
//OpenTelemetry
const opentelemetry = require('@opentelemetry/sdk-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');
const {Resource} = require('@opentelemetry/resources');
const {SemanticResourceAttributes} = require('@opentelemetry/semantic-conventions');

const exporterOptions = {
    url: 'http://localhost:4318/v1/traces'
  }

const traceExporter = new OTLPTraceExporter(exporterOptions);
const sdk = new opentelemetry.NodeSDK({
  traceExporter,
  instrumentations: [getNodeAutoInstrumentations()],
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'sampleNestjsApplication'
  })
  });
  
  // initialize the SDK and register with the OpenTelemetry API
  // this enables the API to record telemetry
  sdk.start()
  
  // gracefully shut down the SDK on process exit
  process.on('SIGTERM', () => {
    sdk.shutdown()
    .then(() => console.log('Tracing terminated'))
    .catch((error) => console.log('Error terminating tracing', error))
    .finally(() => process.exit(0));
    });
    
  module.exports = sdk
```

Step 3. Import the tracer module where your app starts
    
```jsx
const tracer = require('./tracer')
```
    

Step 4. Start the tracer

In the `async function boostrap` section of the application code, initialize the tracer as follows: 

```jsx
const tracer = require('./tracer')

import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
  // All of your application code and any imports that should leverage
  // OpenTelemetry automatic instrumentation must go here.

async function bootstrap() {
    await tracer.start();
    const app = await NestFactory.create(AppModule);
    await app.listen(3001);
  }
  bootstrap();
```

Step 5. Run the application