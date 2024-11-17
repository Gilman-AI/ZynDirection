# ZynDirection: insanely lightweight URL shortener

<div align="center">
 
![zyndirection](https://github.com/user-attachments/assets/9c16e9ce-f75f-4098-8b08-024da5d3d448)

</div>

### **ZynDirection** is a single-script, pure-Golang HTTP server with one purpose: **serve URL redirects on Google Cloud Run** based off JSON directives in a Cloud Storage bucket.

Using only 128MB of RAM and a quarter of a vCPU.

**ZynDirection** gets its name from `https://zyn.wtf/<ALIAS-CODE>`, the URL endpoint for which it provides backend service. There is currently no public-facing content on that domain, however.

## How ZynDirection?

I open-sourced this mainly because it's so mind-bogglingly lightweight—smaller than anything I could find on StackOverflow, or even Golang tutorials.

### _**ZynDirection**'s code contains literally nothing but the application logic_. It does this by taking advantage of:

* Google Cloud Run, [a “serverless” computing environment based on containerized code](https://cloud.google.com/run/docs/overview/what-is-cloud-run).

  * Cloud Run automatically handles:

    * [autoscaling from 0 to hundreds of instances](https://cloud.google.com/run/docs/about-instance-autoscaling)
    * [load balancing/traffic splitting](https://cloud.google.com/run/docs/rollouts-rollbacks-traffic-migration)
    * [HTTPS certificate management](https://cloud.google.com/run/docs/mapping-custom-domains)
    * [TLS negotiation _a.k.a._ reverse proxying](https://cloud.google.com/run/docs/triggering/https-request), and
    * [DDoS protection](https://cloud.google.com/run/docs/securing/security).

* Google Cloud Storage, another “serverless” GCP product which offers [a truly global, massively-replicated storage network](https://cloud.google.com/storage/docs/locations).

  * Unlike almost all other cloud systems, [GCS only bills for storage—not compute](https://cloud.google.com/storage/pricing).

  * For an application like **ZynDirection**, where each user makes *at most* one request and content is not being continuously served, the latency of a GCS bucket is competetive with a high-speed, in-memory data store like Redis.

  * This is because [Cloud Run nodes are either colocated with, or located on the same fiber-optic internal GCP network as, the Colossus filesystem clusters behind Cloud Storage](https://cloud.google.com/storage/docs/google-integration).

* [The Docker `scratch` template image](https://hub.docker.com/_/scratch), which contains **literally nothing**!

  * You will notice in our [Dockerfile](src/Dockerfile) that we make use of [multi-stage Docker builds](https://docs.docker.com/build/building/multi-stage/).

  * Initially, we load `golang:1.23.3-bookworm`, the [official Golang image which is ~800MB and way too large](https://hub.docker.com/_/golang) to load on a dime if we are autoscaling from 0 instances.

  * After executing the Golang build toolchain and compiling (optimized) binaries, we then copy these into a separate Docker container, `scratch`, which is even leaner than the famous `alpine` Linux distro—itself less than 20MB.

  * **ZynDirection**’s Docker image clocks in at 8MB! That’s, like, the size of a iPhone photograph, not a containerized runtime.

The **ZynDirection** container _loads in <1s,_ so users experience no latency [even during so-called “cold starts.”](https://cloud.google.com/functions/docs/bestpractices/tips#min)

We also make use of [Cloud Build](https://cloud.google.com/build/docs), compiling **ZynDirection** by running the following command from inside the `src/` directory:

```bash
gcloud builds submit --async --tag <GCP-REGION>-docker.pkg.dev/<PROJECT-NAME>/zyndirection/zdserver:1.0
```

## Authors

**ZynDirection** was originally built as an internal tool for HFG Ventures LLC. It has no third-party Golang dependencies, except for the [Google Cloud Client Libraries for Go](https://cloud.google.com/go/docs/reference).

## License

**ZynDirection** is licensed under the Apache License, Version 2.0. You may not use the files in this repository except in compliance with the License. You may obtain a copy of the License at [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0), and a copy is included as [the LICENSE file in the root directory of this repository](LICENSE). Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

---
Copyright © 2024 Herbert F. Gilman _a.k.a._ HFG Ventures LLC. See [the NOTICE file](NOTICE.md) for further copyright and licensing information.
