import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useUsers } from '@/composables/useUsers'
import type { User } from '@/types/User'

const sampleUser: User = { id: 1, name: 'Ada' }

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn())
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('useUsers', () => {
  it('starts with empty state', () => {
    const { users, error, isLoading } = useUsers()
    expect(users.value).toEqual([])
    expect(error.value).toBeNull()
    expect(isLoading.value).toBe(false)
  })

  it('fetchUsers hydrates users from the API response', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => [sampleUser],
    } as Response)

    const { users, error, fetchUsers } = useUsers()
    await fetchUsers()

    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/users')
    expect(users.value).toEqual([sampleUser])
    expect(error.value).toBeNull()
  })

  it('fetchUsers records an error when the request fails', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      statusText: 'Internal Server Error',
    } as Response)

    const { users, error, fetchUsers } = useUsers()
    await fetchUsers()

    expect(users.value).toEqual([])
    expect(error.value).toContain('Internal Server Error')
  })

  it('createUser posts the new user and refreshes the list', async () => {
    vi.mocked(fetch)
      .mockResolvedValueOnce({ ok: true } as Response) // POST /users
      .mockResolvedValueOnce({ ok: true, json: async () => [sampleUser] } as Response) // refresh GET

    const { users, createUser } = useUsers()
    await createUser({ name: 'Ada' })

    expect(fetch).toHaveBeenNthCalledWith(
      1,
      'http://localhost:8080/users',
      expect.objectContaining({ method: 'POST' }),
    )
    expect(fetch).toHaveBeenNthCalledWith(2, 'http://localhost:8080/users')
    expect(users.value).toEqual([sampleUser])
  })
})
