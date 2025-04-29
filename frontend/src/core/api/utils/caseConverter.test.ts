import { snakeToCamelCase } from './caseConverter';

describe('snakeToCamelCase', () => {
  it('should convert snake_case keys to camelCase', () => {
    const input = {
      user_id: 1,
      first_name: 'John',
      last_name: 'Doe'
    };

    const expected = {
      userId: 1,
      firstName: 'John',
      lastName: 'Doe'
    };

    expect(snakeToCamelCase(input)).toEqual(expected);
  });

  it('should handle nested objects', () => {
    const input = {
      user_details: {
        home_address: {
          street_name: 'Main St',
          zip_code: '12345'
        }
      }
    };

    const expected = {
      userDetails: {
        homeAddress: {
          streetName: 'Main St',
          zipCode: '12345'
        }
      }
    };

    expect(snakeToCamelCase(input)).toEqual(expected);
  });

  it('should handle arrays', () => {
    const input = {
      user_ids: [1, 2, 3],
      addresses: [
        { address_type: 'home', street_name: 'First St' },
        { address_type: 'work', street_name: 'Second St' }
      ]
    };

    const expected = {
      userIds: [1, 2, 3],
      addresses: [
        { addressType: 'home', streetName: 'First St' },
        { addressType: 'work', streetName: 'Second St' }
      ]
    };

    expect(snakeToCamelCase(input)).toEqual(expected);
  });

  it('should not modify primitive values', () => {
    expect(snakeToCamelCase(null)).toBeNull();
    expect(snakeToCamelCase(undefined)).toBeUndefined();
    expect(snakeToCamelCase(123)).toBe(123);
    expect(snakeToCamelCase('hello')).toBe('hello');
    expect(snakeToCamelCase(true)).toBe(true);
  });

  it('should handle empty objects and arrays', () => {
    expect(snakeToCamelCase({})).toEqual({});
    expect(snakeToCamelCase([])).toEqual([]);
  });
});