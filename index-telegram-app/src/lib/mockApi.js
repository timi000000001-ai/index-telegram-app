/**
 * 生成模拟搜索结果数据
 * @param {string} query - 搜索查询
 * @param {number} page - 当前页码
 * @param {number} size - 每页大小
 * @returns {{results: any[], total: number}}
 */
export function generateMockData(query, page, size) {
  const totalResults = 100;
  const results = [];
  const start = (page - 1) * size;
  const end = start + size;

  for (let i = start; i < end && i < totalResults; i++) {
    results.push({
      id: i + 1,
      title: `关于“${query}”的模拟结果 ${i + 1}`,
      content: `这是关于“${query}”的第 ${i + 1} 条模拟搜索结果的详细内容。此内容由模拟数据生成器提供，用于在API不可用时提供前端页面测试和展示。`,
      source: `模拟来源 ${i % 5 + 1}`,
      type: ['group', 'channel', 'private', 'media', 'link'][i % 5],
      timestamp: new Date(Date.now() - Math.random() * 1000000000).toISOString(),
      relevance: Math.random(),
      interactions: Math.floor(Math.random() * 1000),
    });
  }

  return {
    results,
    total: totalResults,
  };
}


/**
 * 生成模拟搜索建议
 * @param {string} q - 搜索查询字符串
 * @returns {string[]} 建议列表
 */
export function generateMockSuggestions(q) {
  const baseSuggestions = [
    'Vue.js开发技巧',
    'React组件设计',
    'JavaScript异步编程',
    'CSS布局方案',
    'Node.js后端开发',
    'TypeScript类型系统',
    'Webpack配置优化',
    'Git版本控制',
    'Docker容器化',
    'API接口设计'
  ];

  return baseSuggestions
    .filter(suggestion => suggestion.toLowerCase().includes(q.toLowerCase()))
    .slice(0, 5);
}


/**
 * 生成模拟热门搜索数据
 * @author 前端开发团队
 * @date 2024-01-22
 * @version 2.0.0
 * @returns {Array<{keyword: string, count: number, category: string, trend: string}>} 热门搜索列表
 * @description 生成包含多种类别的热门搜索数据，支持趋势标识和分类展示
 */
export function generateMockTrending() {
  // 定义不同类别的热门搜索关键词
  const trendingCategories = {
    technology: [

      { keyword: 'Vue3 Composition API', count: 2340, trend: 'up' },
      { keyword: 'React Server Components', count: 1980, trend: 'hot' },
      { keyword: 'TypeScript 5.0新特性', count: 1756, trend: 'up' },
      { keyword: 'Vite 4.0构建优化', count: 1542, trend: 'stable' },
      { keyword: 'Next.js 13 App Router', count: 1423, trend: 'up' }
    ],
    blockchain: [
      { keyword: '以太坊2.0升级', count: 1890, trend: 'hot' },
      { keyword: 'Web3开发入门', count: 1654, trend: 'up' },
      { keyword: 'NFT智能合约', count: 1234, trend: 'stable' },
      { keyword: 'DeFi协议分析', count: 987, trend: 'down' }
    ],
    ai: [
      { keyword: '机器学习算法', count: 2156, trend: 'hot' },
      { keyword: '深度学习框架', count: 1876, trend: 'up' },
      { keyword: '自然语言处理', count: 1543, trend: 'stable' },
      { keyword: '计算机视觉', count: 1321, trend: 'up' }
    ],
    mobile: [
      { keyword: 'Flutter跨平台开发', count: 1678, trend: 'up' },
      { keyword: 'React Native性能优化', count: 1456, trend: 'stable' },
      { keyword: 'iOS SwiftUI', count: 1234, trend: 'up' },
      { keyword: 'Android Jetpack Compose', count: 1123, trend: 'hot' }
    ],
    devops: [
      { keyword: 'Docker容器化部署', count: 1789, trend: 'stable' },
      { keyword: 'Kubernetes集群管理', count: 1567, trend: 'up' },
      { keyword: 'CI/CD自动化流程', count: 1345, trend: 'stable' },
      { keyword: '微服务架构设计', count: 1234, trend: 'up' }
    ]
  };

  // 合并所有类别的数据
  /** @type {Array<{keyword: string, count: number, trend: string, category: string}>} */
  const allTrending = [];
  Object.entries(trendingCategories).forEach(([category, items]) => {
    items.forEach(item => {
      allTrending.push({
        ...item,
        category: category
      });
    });
  });



  // 按搜索次数排序并返回前20个
  return allTrending
    .sort((a, b) => b.count - a.count)
    .slice(0, 20)
    .map((item, index) => ({
      ...item,
      rank: index + 1,
      // 添加一些随机波动使数据更真实
      count: item.count + Math.floor(Math.random() * 100) - 50
    }));
}