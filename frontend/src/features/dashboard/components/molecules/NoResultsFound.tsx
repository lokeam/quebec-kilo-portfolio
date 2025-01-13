export function NoResultsFound() {
  return (
    <div className="text-center py-12">
      <h3 className="text-lg font-medium">No results found</h3>
      <p className="text-sm text-gray-500">
        Hmm. We couldn't find any results for your search.
      </p>
      <p className="text-sm text-gray-500">
        Try adjusting your filters or search terms.
      </p>
    </div>
  );
}