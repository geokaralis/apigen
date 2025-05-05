/**
 * This source file is auto-generated.
 */

export type Memory = {
  /** Total memory size in MB */
  total_size: number;
  /** Memory type */
  type: string;
  /** Memory speed in MHz */
  speed?: number;
};

export type Storage = {
  /** Total storage size in GB */
  total_size: number;
  /** Storage type */
  type: "ssd" | "hdd" | "nvme";
  /** IOPS (Input/Output Operations Per Second) */
  iops?: number;
};

export type Network = {
  /** Whether a public IP is assigned */
  public_ip?: boolean;
  /** Maximum Transmission Unit (MTU) */
  mtu?: number;
  /** Network bandwidth in Mbps */
  bandwidth: number;
};

export type Gpu = {
  /** GPU model name */
  model: string;
  /** GPU memory in MB */
  memory: number;
  /** Number of CUDA cores, if applicable */
  cuda_cores?: number;
};

export type Error = {
  /** Error code */
  code: string;
  /** Error message */
  message: string;
};

export type ResourceResultsPage = {
  /** List of Resources */
  items?: Resource[];
  /** Total number of items available */
  total?: number;
  /** Number of items per page */
  limit?: number;
  /** Current offset (starting index) */
  offset?: number;
};

export type Resource = {
  /** Region where the resource is located */
  region?: string;
  /** Availability zone where the resource is located */
  zone?: string;
  /** Resource-specific specifications based on type */
  specifications?: Record<string, any>;
  /** Unique identifier of the resource */
  id: string;
  /** User-friendly name of the resource */
  name: string;
  /** Current status of the resource */
  status: "available" | "in_use" | "maintenance" | "error";
  /** Current utilization as a percentage */
  utilization_percentage?: number;
  /** User-defined tags for the resource */
  tags?: Record<string, any>;
  /** Creation timestamp */
  created_at?: string;
  /** Last update timestamp */
  updated_at?: string;
  /** Type of computing resource */
  type: "cpu" | "memory" | "storage" | "network" | "gpu";
  /** ID of the VPC this resource belongs to */
  vpc_id?: string;
};

export type Cpu = {
  /** Number of CPU cores */
  cores: number;
  /** CPU architecture */
  architecture: string;
  /** Base clock speed in GHz */
  clock_speed?: number;
  /** Number of threads per core */
  threads_per_core?: number;
};

/**
 * Parameters for resource_list
 * Returns a list of computing resources with optional filtering
 */
export interface ResourceListParams {
  /** Filter resources by VPC name or ID */
  vpc?: string;
  /** Filter resources by type */
  type?: "cpu" | "memory" | "storage" | "network" | "gpu";
  /** Filter resources by status */
  status?: "available" | "in_use" | "maintenance" | "error";
  /** Filter resources by region */
  region?: string;
  /** Maximum number of items to return */
  limit?: number;
  /** Number of items to skip */
  offset?: number;
}

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  /**
   * List computing resources
   * Returns a list of computing resources with optional filtering
   */
  async resourceList(
    params?: ResourceListParams
  ): Promise<ResourceResultsPage> {
    const queryParams = new URLSearchParams();
    if (params) {
      if (params.vpc !== undefined) {
        queryParams.append("vpc", String(params.vpc));
      }
      if (params.type !== undefined) {
        queryParams.append("type", String(params.type));
      }
      if (params.status !== undefined) {
        queryParams.append("status", String(params.status));
      }
      if (params.region !== undefined) {
        queryParams.append("region", String(params.region));
      }
      if (params.limit !== undefined) {
        queryParams.append("limit", String(params.limit));
      }
      if (params.offset !== undefined) {
        queryParams.append("offset", String(params.offset));
      }
    }

    const url = `${this.baseUrl}/resources${
      queryParams.toString() ? "?" + queryParams.toString() : ""
    }`;
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Request failed");
    }

    return response.json();
  }
  /**
   * Fetch a single resource
   * Returns detailed information about a specific computing resource
   */
  async resourceView(resource: string): Promise<Resource> {
    const url = `${this.baseUrl}/resources/${resource}`;
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Request failed");
    }

    return response.json();
  }
}
