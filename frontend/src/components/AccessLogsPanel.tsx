import { useEffect, useMemo, useState } from 'react'
import { fetchAccessLogs } from '../api/client'
import type { AccessLogEntry } from '../types'

export default function AccessLogsPanel() {
  const [logs, setLogs] = useState<AccessLogEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [autoScroll, setAutoScroll] = useState(false)

  const load = async () => {
    try {
      setLoading(true)
      setError('')
      const res = await fetchAccessLogs()
      setLogs(res.logs.slice(-200))
    } catch (e) {
      setError(e instanceof Error ? e.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void load() }, [])
  useEffect(() => {
    if (!autoScroll) return
    const id = window.setInterval(() => { void load() }, 3000)
    return () => window.clearInterval(id)
  }, [autoScroll])

  useEffect(() => {
    if (autoScroll) window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' })
  }, [logs, autoScroll])

  const items = useMemo(() => [...logs].reverse(), [logs])

  return (
    <div className="p-4 md:p-6 space-y-4">
      <div className="flex flex-wrap items-center gap-3 justify-between">
        <div>
          <h2 className="text-xl font-bold">运行日志</h2>
          <p className="text-sm text-base-content/60">代理访问日志，保留最近 200 条</p>
        </div>
        <div className="flex items-center gap-3">
          <label className="label cursor-pointer gap-2">
            <span className="label-text">自动滚动</span>
            <input type="checkbox" className="toggle toggle-primary" checked={autoScroll} onChange={(e) => setAutoScroll(e.target.checked)} />
          </label>
          <button className="btn btn-primary btn-sm" onClick={() => void load()} disabled={loading}>刷新</button>
        </div>
      </div>

      {error && <div className="alert alert-error"><span>{error}</span></div>}
      {loading && <div className="text-sm text-base-content/60">加载中...</div>}

      <div className="space-y-3">
        {items.map((log, idx) => (
          <div key={`${log.timestamp}-${idx}`} className="card bg-base-100 border border-base-300 shadow-sm">
            <div className="card-body p-4 gap-2">
              <div className="flex flex-wrap items-center gap-2 text-sm">
                <span className="badge badge-outline">{log.status === 'success' ? '成功' : '失败'}</span>
                <span className="font-mono text-xs opacity-70">{log.timestamp}</span>
              </div>
              <div className="grid md:grid-cols-2 gap-2 text-sm">
                <div><span className="opacity-60">来源 IP：</span><span className="font-mono">{log.source_ip || '-'}</span></div>
                <div><span className="opacity-60">目标地址：</span><span className="font-mono break-all">{log.target || '-'}</span></div>
                <div><span className="opacity-60">出站节点：</span><span>{log.outbound_node || '-'}</span></div>
                <div><span className="opacity-60">结果：</span><span>{log.status === 'success' ? '成功' : `失败${log.error ? `：${log.error}` : ''}`}</span></div>
              </div>
            </div>
          </div>
        ))}
        {!loading && items.length === 0 && <div className="text-sm text-base-content/60">暂无访问日志</div>}
      </div>
    </div>
  )
}
