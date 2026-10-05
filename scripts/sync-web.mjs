import { cp, mkdir, readdir, rm } from "node:fs/promises";
import { resolve } from "node:path";

const source = resolve("apps/web/out");
const target = resolve("cmd/agent/web");

await mkdir(target, { recursive: true });
for (const entry of await readdir(target)) {
  if (entry !== "placeholder.html") {
    await rm(resolve(target, entry), { recursive: true, force: true });
  }
}
await cp(source, target, { recursive: true });
