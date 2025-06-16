-- Clean up in reverse order of creation
DELETE FROM yearly_spending_aggregates WHERE user_id = '9a4aeee6-fb31-4839-a921-f61b0525046d';
DELETE FROM monthly_spending_aggregates WHERE user_id = '9a4aeee6-fb31-4839-a921-f61b0525046d';
DELETE FROM digital_location_subscriptions WHERE digital_location_id IN (
    '9a4aeee6-fb31-4839-a921-f61b0525046d',
    '9a4aeee6-fb31-4839-a921-f61b0525046e',
    '9a4aeee6-fb31-4839-a921-f61b0525046f',
    '9a4aeee6-fb31-4839-a921-f61b05250470'
);
DELETE FROM digital_locations WHERE id IN (
    '9a4aeee6-fb31-4839-a921-f61b0525046d',
    '9a4aeee6-fb31-4839-a921-f61b0525046e',
    '9a4aeee6-fb31-4839-a921-f61b0525046f',
    '9a4aeee6-fb31-4839-a921-f61b05250470',
    '9a4aeee6-fb31-4839-a921-f61b05250471'
);
DELETE FROM one_time_purchases WHERE user_id = '9a4aeee6-fb31-4839-a921-f61b0525046d';
DELETE FROM spending_categories WHERE name IN ('hardware', 'dlc', 'in_game', 'subscription', 'physical', 'disc');
DELETE FROM users WHERE id = '9a4aeee6-fb31-4839-a921-f61b0525046d';