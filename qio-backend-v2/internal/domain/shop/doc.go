// Package shop 承载商城域：信纸、字体、字体颜色、印章、功能卡的展示与购买，
// 以及运营位的外部文创商品。
//
// v1 中这几类道具各自一个模块（FontService、PaperService、CardService、
// MarketingService），业务形态一致——列表展示、购买扣费、写入背包——在 v2 合并
// 为同一个域，用 ItemType 区分。
//
// 本域是纯读的。购买动作的扣费与背包写入都发生在 user 表，因此「购买」这一用例
// 需要同时操作 shop 与 user 两个域，编排放在 internal/app 而不是本域的 service。
// 使用功能卡还会改动信件，涉及三个域，同样归 internal/app。
//
// 印章（signet）在 v1 只作为注册赠品，没有购买接口，因此本域不提供其查询方法。
package shop
