require 'json'

Class.new(Nanoc::Filter) do
  identifiers :to_search_json

  def run(content, arguments={})
    JSON.dump(
      title: @item[:title] || @item[:short_title],
      content: Nokogiri::HTML.fragment(content).inner_text,
      path: @item.path,
    )
  end
end
