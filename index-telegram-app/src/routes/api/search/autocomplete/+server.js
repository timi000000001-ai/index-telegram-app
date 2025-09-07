import { MeiliSearch } from 'meilisearch'
import { json } from '@sveltejs/kit'

const client = new MeiliSearch({
  host: 'http://127.0.0.1:7700',
  apiKey: 'timigogogo',
})

/** @type {import('./$types').RequestHandler} */
export async function GET({ url }) {
  const query = url.searchParams.get('q')

  if (!query) {
    return json({ error: 'Query parameter \"q\" is required' }, { status: 400 })
  }

  try {
    const searchResult = await client.index('suggestions').search(query, {
      limit: 8, // Limit the number of suggestions
    })

    // We only need the 'query' field from the results
    const suggestions = searchResult.hits.map(item => item.query)

    return json(suggestions)
  } catch (error) {
    console.error('Meilisearch error:', error)
    return json({ error: 'Failed to fetch suggestions' }, { status: 500 })
  }
}