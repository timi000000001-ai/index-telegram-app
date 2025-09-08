import MeiliSearch from 'meilisearch';
import { json } from '@sveltejs/kit';

const client = new MeiliSearch({
  host: 'http://localhost:7700',
  apiKey: 'masterKey',
});

export async function GET({ url }) {
  const query = url.searchParams.get('q');

  if (!query) {
    return json({ error: 'Query parameter "q" is required' }, { status: 400 });
  }

  try {
    const searchResult = await client.index('suggestions').search(query, {
      limit: 10,
    });
    return json(searchResult.hits);
  } catch (error) {
    console.error('Meilisearch error:', error);
    return json({ error: 'Failed to fetch suggestions' }, { status: 500 });
  }
}