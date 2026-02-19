export const PORT =
    import.meta.env.VITE_PORT ?? 3003;

export const API_BASE =
    import.meta.env.VITE_API_BASE ?? `http://10.88.81.1:${PORT}/api`;