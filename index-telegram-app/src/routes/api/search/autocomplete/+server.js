import MeiliSearch from 'meilisearch';
import { json } from '@sveltejs/kit';

const client = new MeiliSearch({
  host: 'http://localhost:7700',
  apiKey: 'timigogogo',
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
    const suggestions = searchResult.hits.map(hit => hit.query);
    return json(suggestions);
  } catch (error) {
    console.error('Meilisearch error:', error);
    return json({ error: 'Failed to fetch suggestions' }, { status: 500 });
  }
}