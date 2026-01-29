-- 初始种子数据

USE yidaiku;

-- 插入默认穿搭模板（用于冷启动）
INSERT INTO outfit_templates (template_id, name, description, items, occasion, style, season, sort_order, is_enabled) VALUES
-- 通勤模板
('tpl_001', '经典通勤', '白衬衫搭配深色西裤，正式又不失优雅', 
 '[{"category":"上衣","description":"白色/浅色衬衫"},{"category":"下装","description":"深色西裤/半身裙"},{"category":"鞋子","description":"黑色/棕色皮鞋"}]', 
 '职场通勤', '正式', '四季通用', 1, 1),

('tpl_002', '商务休闲', '针织衫搭配直筒裤，舒适又得体', 
 '[{"category":"上衣","description":"纯色针织衫"},{"category":"下装","description":"直筒裤"},{"category":"鞋子","description":"乐福鞋/平底鞋"}]', 
 '职场通勤', '通勤', '四季通用', 2, 1),

-- 休闲模板
('tpl_003', '休闲日常', '经典白T配牛仔裤，简约百搭', 
 '[{"category":"上衣","description":"白T/条纹T"},{"category":"下装","description":"牛仔裤/休闲裤"},{"category":"鞋子","description":"小白鞋/运动鞋"}]', 
 '周末逛街', '休闲', '四季通用', 3, 1),

('tpl_004', '街头时尚', '卫衣搭配工装裤，潮流感十足', 
 '[{"category":"上衣","description":"纯色/印花卫衣"},{"category":"下装","description":"工装裤"},{"category":"鞋子","description":"老爹鞋/板鞋"}]', 
 '周末逛街', '休闲', '秋冬', 4, 1),

-- 约会模板
('tpl_005', '甜美约会', '碎花裙搭配小白鞋，清新甜美', 
 '[{"category":"上衣","description":"碎花连衣裙"},{"category":"鞋子","description":"小白鞋/凉鞋"},{"category":"配饰","description":"小包/发饰"}]', 
 '约会', '甜美', '春夏', 5, 1),

('tpl_006', '优雅约会', '针织上衣搭配半身裙，优雅大方', 
 '[{"category":"上衣","description":"纯色针织上衣"},{"category":"下装","description":"A字半身裙"},{"category":"鞋子","description":"高跟鞋/尖头平底鞋"}]', 
 '约会', '正式', '四季通用', 6, 1),

-- 运动模板
('tpl_007', '日常运动', '运动T恤搭配运动裤，舒适自在', 
 '[{"category":"上衣","description":"运动T恤/速干衣"},{"category":"下装","description":"运动裤/运动短裤"},{"category":"鞋子","description":"运动鞋"}]', 
 '运动', '运动', '四季通用', 7, 1),

-- 面试模板
('tpl_008', '正式面试', '西装外套搭配衬衫，专业有型', 
 '[{"category":"上衣","description":"深色西装外套"},{"category":"内搭","description":"白色衬衫"},{"category":"下装","description":"西裤/西装裙"},{"category":"鞋子","description":"黑色皮鞋/高跟鞋"}]', 
 '面试', '正式', '四季通用', 8, 1);

-- 插入默认搭配规则
INSERT INTO matching_rules (rule_id, rule_type, name, `condition`, action, weight, priority, is_enabled) VALUES
-- 颜色搭配规则
('rule_color_001', 'color', '中性色百搭', 
 '{"top_color":["黑色","白色","灰色","米色"]}', 
 '{"bottom_colors":["any"]}', 
 0.95, 100, 1),

('rule_color_002', 'color', '同色系搭配', 
 '{"same_color_family":true}', 
 '{"score_boost":0.1}', 
 0.90, 90, 1),

('rule_color_003', 'color', '深蓝色搭配', 
 '{"top_color":["深蓝色"]}', 
 '{"bottom_colors":["白色","黑色","米色","卡其色"]}', 
 0.85, 80, 1),

('rule_color_004', 'color', '避免红绿搭配', 
 '{"top_color":["红色"],"bottom_color":["绿色"]}', 
 '{"score_penalty":-0.3}', 
 0.70, 70, 1),

-- 场合适配规则
('rule_occ_001', 'occasion', '职场通勤', 
 '{"occasion":"职场通勤"}', 
 '{"styles":["正式","通勤"],"colors":["深色系","中性色"],"avoid":["运动鞋","破洞牛仔裤"]}', 
 0.90, 100, 1),

('rule_occ_002', 'occasion', '约会场合', 
 '{"occasion":"约会"}', 
 '{"styles":["甜美","优雅","时尚"],"colors":["亮色","柔和色系"],"prefer":["裙子","有设计感的上衣"]}', 
 0.90, 100, 1),

('rule_occ_003', 'occasion', '运动场合', 
 '{"occasion":"运动"}', 
 '{"styles":["运动"],"prefer":["运动鞋","运动裤"]}', 
 0.95, 100, 1),

('rule_occ_004', 'occasion', '面试场合', 
 '{"occasion":"面试"}', 
 '{"styles":["正式"],"colors":["黑色","深蓝色","深灰色"],"avoid":["鲜艳颜色","休闲单品"]}', 
 0.95, 100, 1),

-- 季节适配规则
('rule_season_001', 'season', '春夏穿搭', 
 '{"season":"春夏","temperature_gt":20}', 
 '{"prefer_season":["春夏","四季通用"],"material":["轻薄","透气"]}', 
 0.90, 100, 1),

('rule_season_002', 'season', '秋冬穿搭', 
 '{"season":"秋冬","temperature_lt":15}', 
 '{"prefer_season":["秋冬","四季通用"],"material":["厚实","保暖"]}', 
 0.90, 100, 1),

-- 风格统一规则
('rule_style_001', 'style', '风格一致性', 
 '{"check_style_consistency":true}', 
 '{"same_style_boost":0.15,"mixed_style_penalty":-0.1}', 
 0.85, 80, 1);
