import * as os from "os";
import * as github from "@actions/github";
import * as fs from "fs";

export function getArch(): string {
  const arch = os.arch();
  if (arch === "x64") {
    return "amd64";
  }
  return arch;
}

export function getDownloadUrl(version: string, arch: string): string {
  switch (os.type()) {
    case "Linux":
      return `https://github.com/1k-off/action-semver-tag-release/releases/download/${version}/semver-tag-release-linux-${arch}`;
    case "Darwin":
      return `https://github.com/1k-off/action-semver-tag-release/releases/download/${version}/semver-tag-release-darwin-${arch}`;
    case "Windows_NT":
      return `https://github.com/1k-off/action-semver-tag-release/releases/download/${version}/semver-tag-release-windows-${arch}.exe`;
    default:
      if (os.type().match(/^Win/)) {
        return `https://github.com/1k-off/action-semver-tag-release/releases/download/${version}/semver-tag-release-windows-${arch}.exe`;
      }
      return "";
  }
}

export function getExecutableExtension(): string {
  if (os.type().match(/^Win/)) {
    return ".exe";
  }
  return "";
}

export function getToolName(): string {
  let arch = getArch();
  switch (os.type()) {
    case "Linux":
      return `semver-tag-release-linux-${arch}`;
    case "Darwin":
      return `semver-tag-release-darwin-${arch}`;
    case "Windows_NT":
      return `semver-tag-release-windows-${arch}.exe`;
    default:
      if (os.type().match(/^Win/)) {
        return `semver-tag-release-windows-${arch}.exe`;
      }
      return "";
  }
}

export async function getLatestVersion(token: string): Promise<string> {
  const octokit = github.getOctokit(token);
  let version = "";
  let response = await octokit.request(
    "GET /repos/1k-off/action-semver-tag-release/releases/latest",
    {
      owner: "1k-off",
      repo: "action-semver-tag-release",
      headers: {
        "X-GitHub-Api-Version": "2022-11-28",
      },
    }
  );
  version = response.data.tag_name;
  if (version == "") {
    throw new Error("Failed to get latest version");
  }
  return version;
}
