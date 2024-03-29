import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage, ElNotification } from 'element-plus'
import { getStorage } from '@/utils/storage'
class Request {
  // axios 实例
  instance: AxiosInstance
  // baseConfig: AxiosRequestConfig = { baseURL: "/api", timeout: 60000 };
  constructor(config: AxiosRequestConfig) {
    axios.defaults.headers.post['Content-Type'] = 'application/json';
    // 使用axios.create创建axios实例
    this.instance = axios.create(config)
    // this.instance = axios.create(Object.assign(this.baseConfig, config));
    this.instance.interceptors.request.use(
      (config: AxiosRequestConfig) => {
        // 一般会请求拦截里面加token
        const token = getStorage('token')
        config.headers!.Authorization = 'Bearer ' + token
        let { data = {}, method } = config
        if (method === 'post' || method === 'put') {
          const formData = new FormData()
          Object.keys(data).forEach((item) => {
            // console.log(data, data[item] instanceof Array, data[item])
            if (data[item] instanceof Array) {
              data[item].forEach((el: any) => {
                formData.append(item + '[]', el)
              })
            } else {
              formData.append(item, data[item])
            }
          })
          config.data = formData
        }
        return config
      },
      (err: any) => {
        return Promise.reject(err)
      }
    )

    this.instance.interceptors.response.use(
      (res: AxiosResponse) => {
        // 直接返回res，当然你也可以只返回res.data
        if (Object.prototype.toString.apply(res.data) === '[object Blob]') {
          return res.data
        }
        const { data, code, message } = res.data
        if (code === 200) {
          return data
        } else {
          ElNotification({
            title: 'Warning',
            message: message,
            type: 'warning',
          })
          return Promise.reject(res)
        }
      },
      (err: any) => {
        ElMessage({
          message: err.message,
          type: 'error',
        })
        // 请求错误
        return Promise.reject(err.response)
      }
    )
  }

  // 定义请求方法
  public request(config: AxiosRequestConfig): Promise<AxiosResponse> {
    return this.instance.request(config)
  }

  public get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.instance.get(url, config)
  }

  public post<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return this.instance.post(url, data, config)
  }

  public put<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    console.log(555)
    return this.instance.put(url, data, config)
  }

  public delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.instance.delete(url, config)
  }
}

export default Request
