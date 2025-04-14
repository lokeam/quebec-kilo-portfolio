# API Standards and Best Practices

## Table of Contents
1. [Axios Response Handling](#axios-response-handling)
2. [Type Definitions](#type-definitions)
3. [Error Handling](#error-handling)
4. [Code Review Checklist](#code-review-checklist)

## Axios Response Handling

### Standard Pattern
```typescript
// CORRECT: Direct data access (use this by default)
const data = await axiosInstance.get<DataType>('/endpoint');

// CORRECT: Full response access (only when needed)
const { data, status, headers } = await axiosInstance.get<DataType>('/endpoint');

// INCORRECT: Don't use AxiosResponse in generic type
const response = await axiosInstance.get<DataType, AxiosResponse<DataType>>('/endpoint'); // ❌
```

### When to Use Each Pattern

1. **Direct Data Access (Default)**
   - Use when you only need the response data
   - Most common use case
   - Simplest and most straightforward

2. **Full Response Access**
   - Use when you need status codes, headers, or other response metadata
   - Example: Checking response status for error handling
   - Example: Accessing response headers for pagination

### Common Mistakes to Avoid

1. **Incorrect Type Usage**
   ```typescript
   // ❌ DON'T: Over-complicate type definitions
   const response = await axiosInstance.get<DataType, AxiosResponse<DataType>>('/endpoint');

   // ✅ DO: Keep it simple
   const data = await axiosInstance.get<DataType>('/endpoint');
   ```

2. **Unnecessary Response Wrapping**
   ```typescript
   // ❌ DON'T: Unnecessarily wrap the response
   const { data } = await axiosInstance.get<DataType>('/endpoint');
   return { data }; // Unnecessary wrapping

   // ✅ DO: Return the data directly
   return await axiosInstance.get<DataType>('/endpoint');
   ```

## Type Definitions

### Keep Types Simple and Purposeful
```typescript
// ✅ DO: Use types for documentation and compile-time checks
interface User {
  id: string;
  name: string;
  email: string;
}

// ❌ DON'T: Over-engineer types
interface ComplexUserResponse<T extends User> {
  data: T;
  metadata: {
    timestamp: string;
    version: string;
  };
  status: number;
}
```

### Type Usage Guidelines

1. **Use Interface for API Responses**
   ```typescript
   // ✅ DO: Simple, clear interface
   interface ApiResponse {
     id: string;
     name: string;
   }
   ```

2. **Avoid Complex Generic Types**
   ```typescript
   // ❌ DON'T: Overly complex generics
   type ApiResponse<T extends BaseType, K extends keyof T> = {
     [P in K]: T[P];
   } & {
     metadata: MetadataType;
   };
   ```

## Error Handling

### Standard Error Handling Pattern
```typescript
try {
  const data = await axiosInstance.get<DataType>('/endpoint');
  return data;
} catch (error) {
  if (axios.isAxiosError(error)) {
    // Handle specific Axios errors
    console.error('API Error:', error.response?.status, error.message);
  } else {
    // Handle other errors
    console.error('Unexpected error:', error);
  }
  return FALLBACK_DATA;
}
```

### Error Handling Guidelines

1. **Always Use Try-Catch**
   - Wrap API calls in try-catch blocks
   - Handle both Axios-specific and general errors

2. **Provide Meaningful Error Messages**
   - Log relevant error details
   - Include status codes and error messages
   - Don't expose sensitive information

## Code Review Checklist

### API Calls
- [ ] Uses correct Axios response pattern
- [ ] Properly typed response data
- [ ] Includes error handling
- [ ] Follows standard patterns documented here

### Type Definitions
- [ ] Types are simple and purposeful
- [ ] No unnecessary complexity
- [ ] Properly documented

### Error Handling
- [ ] Includes try-catch blocks
- [ ] Handles both Axios and general errors
- [ ] Provides meaningful error messages
- [ ] Includes appropriate fallback behavior

## Examples

### Good Example
```typescript
// Simple, clear, and follows standards
const getUser = async (id: string): Promise<User> => {
  try {
    const user = await axiosInstance.get<User>(`/users/${id}`);
    return user;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error('Failed to fetch user:', error.message);
    }
    throw error;
  }
};
```

### Bad Example
```typescript
// Overly complex and violates standards
const getUser = async <T extends User>(
  id: string,
  options?: UserOptions
): Promise<AxiosResponse<T>> => {
  const response = await axiosInstance.get<T, AxiosResponse<T>>(
    `/users/${id}`,
    { params: options }
  );
  return {
    data: response.data,
    status: response.status,
    headers: response.headers
  };
};
```