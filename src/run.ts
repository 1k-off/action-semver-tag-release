import * as path from "path";
import * as util from "util";
import * as fs from "fs";

import * as toolCache from "@actions/tool-cache";
import * as core from "@actions/core";
import * as exec from "@actions/exec";

import {
  getDownloadUrl,
  getArch,
  getToolName,
  getExecutableExtension,
  getLatestVersion,
} from "./helpers";

let toolName = getToolName();


export async function run() {
    try {
      let token = core.getInput("github_token", { required: true });
      let version = core.getInput("version", { required: true });

      let repo = process.env.GITHUB_REPOSITORY;
      let repoParts = repo.split("/");
      let owner = repoParts[0];
      let repoName = repoParts[1];
      let sha = process.env.GITHUB_SHA;
      let releaseTag = core.getInput("release_tag", { required: false });
      let preReleaseTag = core.getInput("pre_release_tag", { required: false });
      let createRelease = core.getInput("create_release", { required: false });
      let assets = core.getInput("assets", { required: false });

      // set env vars
      core.exportVariable("GITHUB_TOKEN", token);
      core.exportVariable("RELEASE_TAG", releaseTag);
      core.exportVariable("PRE_RELEASE_TAG", preReleaseTag);
      core.exportVariable("CREATE_RELEASE", createRelease);
      core.exportVariable("ASSETS", assets);

      if (version.toLocaleLowerCase() === "latest") {
        version = await getLatestVersion(token);
      }
      const cachedPath = await downloadTool(version);
      core.addPath(path.dirname(cachedPath));

      core.debug(
          `semver-tag-release version: '${version}' has been cached at ${cachedPath}`
      );
      core.setOutput("semver-tag-release-path", cachedPath);
      // run semver-tag-release
      await exec.exec(cachedPath, [
        '--repo-owner', owner,
        '--repo-name', repoName,
        '--commit-sha', sha,
        '--github-token', token,
        '--release-tag', releaseTag,
        '--pre-release-tag', preReleaseTag,
        '--create-release', createRelease,
        '--assets', assets
      ]);
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
}

export async function downloadTool(version: string): Promise<string> {
  let cachedToolpath = toolCache.find(toolName, version);
  let downloadPath = "";
  const arch = getArch();
  if (!cachedToolpath) {
    try {
      downloadPath = await toolCache.downloadTool(
        getDownloadUrl(version, arch)
      );
    } catch (exception) {
      if (
        exception instanceof toolCache.HTTPError &&
        exception.httpStatusCode === 404
      ) {
        throw new Error(
          util.format(
            "semver-tag-release '%s' for '%s' arch not found.",
            version,
            arch
          )
        );
      } else {
        throw new Error("Download failed");
      }
    }

    cachedToolpath = await toolCache.cacheFile(
      downloadPath,
      toolName + getExecutableExtension(),
      toolName,
      version
    );
  }

  const toolPath = path.join(
    cachedToolpath,
    toolName + getExecutableExtension()
  );
  fs.chmodSync(toolPath, "775");
  return toolPath;
}

run().catch(core.setFailed);
