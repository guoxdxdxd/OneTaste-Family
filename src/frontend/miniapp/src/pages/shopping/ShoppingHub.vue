<template>
  <div class="shopping-page">
    <section class="card hero">
      <div>
        <p class="eyebrow">买菜 · 清单中心</p>
        <h1>把菜单转成行动，避免多买或漏买</h1>
        <p class="subtitle">
          根据需求文档预留生成、编辑、分享清单的区域，未来可直接串联菜单模块与库存管理。
        </p>
      </div>
      <ul class="highlights">
        <li>按菜单自动汇总食材</li>
        <li>分类展示 · 一目了然</li>
        <li>存放时间提醒</li>
      </ul>
    </section>

    <section class="generator card">
      <header>
        <div>
          <h2>自动生成购物清单</h2>
          <p>选择日期范围或直接引用菜单，系统自动合并相同食材并计算总量。</p>
        </div>
        <button type="button">即将接入</button>
      </header>
      <form class="generator-form">
        <label>
          <span>开始日期</span>
          <input type="date" disabled />
        </label>
        <label>
          <span>结束日期</span>
          <input type="date" disabled />
        </label>
        <label>
          <span>关联菜单</span>
          <select disabled>
            <option>本周菜单（3份）</option>
          </select>
        </label>
      </form>
    </section>

    <section class="list-preview card">
      <header>
        <h2>清单内容展示</h2>
        <p>预留按分类展示、数量、单位、存放天数与状态切换的区域。</p>
      </header>
      <div class="categories">
        <article v-for="category in categories" :key="category.name" class="category-card">
          <header>
            <h3>{{ category.name }}</h3>
            <span>{{ category.count }} 项</span>
          </header>
          <ul>
            <li v-for="item in category.items" :key="item.name">
              <div>
                <strong>{{ item.name }}</strong>
                <p>{{ item.quantity }}</p>
              </div>
              <button type="button">标记</button>
            </li>
          </ul>
        </article>
      </div>
    </section>

    <section class="card editor">
      <div>
        <h2>手动编辑清单</h2>
        <p>添加/删除/修改数量、标记已购买等交互入口在此保留，方便后续直接绑定数据。</p>
      </div>
      <div class="editor-actions">
        <button type="button">添加食材</button>
        <button type="button">导入库存</button>
        <button type="button">分享清单</button>
      </div>
    </section>
  </div>
</template>

<script setup>
const categories = [
  {
    name: '蔬菜类',
    count: 4,
    items: [
      { name: '西兰花', quantity: '2 颗 / 2 天内食用' },
      { name: '生菜', quantity: '1 颗 / 3 天内食用' }
    ]
  },
  {
    name: '肉类',
    count: 3,
    items: [
      { name: '鸡胸肉', quantity: '600g / 可冷冻 3 天' },
      { name: '五花肉', quantity: '500g / 今日烹饪' }
    ]
  },
  {
    name: '调料',
    count: 2,
    items: [
      { name: '生抽', quantity: '200ml' },
      { name: '蚝油', quantity: '100ml' }
    ]
  }
]
</script>

<style scoped>
.shopping-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 20px 120px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.card {
  background: var(--color-card);
  border-radius: var(--radius-large);
  padding: 28px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-card);
}

.hero {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: uppercase;
  color: var(--color-text-secondary);
}

.subtitle {
  margin: 0;
  color: var(--color-text-secondary);
}

.highlights {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin: 0;
  padding: 0;
  list-style: none;
}

.highlights li {
  background: var(--color-surface);
  padding: 10px 16px;
  border-radius: var(--radius-medium);
}

.generator header {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.generator button {
  border: none;
  padding: 10px 18px;
  border-radius: var(--radius-medium);
  background: linear-gradient(120deg, var(--color-accent), var(--color-accent-soft));
  color: #fff;
  cursor: pointer;
}

.generator-form {
  margin-top: 20px;
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

.generator-form label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  color: var(--color-text-secondary);
}

.generator-form input,
.generator-form select {
  border-radius: var(--radius-medium);
  border: 1px solid var(--color-border);
  padding: 12px;
  background: var(--color-surface);
}

.list-preview header {
  margin-bottom: 16px;
}

.categories {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.category-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-medium);
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.category-card header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.category-card ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.category-card li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-radius: var(--radius-small);
  background: var(--color-surface);
}

.category-card strong {
  display: block;
}

.category-card p {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.category-card button,
.editor-actions button {
  border: none;
  border-radius: var(--radius-medium);
  padding: 8px 14px;
  background: var(--color-card);
  border: 1px solid var(--color-border);
  cursor: pointer;
}

.editor {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.editor-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.editor-actions button {
  background: var(--color-surface);
}

@media (min-width: 720px) {
  .hero {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }

  .generator header {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
