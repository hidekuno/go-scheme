/*
   Go lang 2nd study program.
   This is c++-11 sample
   hidekuno@gmail.com

   g++ --std=c++11 tree.cpp -o tree
*/
#include <list>
#include <string>
#include <iostream>
#include <memory>
#include <unordered_map>
#include <vector>
#include <sstream>
#include <iostream>
#include <cstdio>
#include <fstream>
#include <stdexcept>

extern "C" {
  #include "unistd.h"
}

using std::string;
using std::list;
using std::cout;
using std::cerr;
using std::endl;
using std::shared_ptr;
using std::weak_ptr;
using std::unordered_map;
using std::vector;
using std::stringstream;
using std::cin;
using std::istream;
using std::ifstream;
using std::exception;
using std::getline;

class Visitor;
class ItemVisitor;
class LineItemVisitor;
class Item;
const char DELIMITER_CHAR = '.';

class Visitor {
public:
	virtual void visit(Item& item) = 0;
};
typedef weak_ptr<Item> ItemPtr;

class Item {
	friend class LineItemVisitor;

private:
	string name;
	shared_ptr<Item> parent;
	list<ItemPtr> children;

public:
	Item(string& s) : name(s) {};
	Item(string& s, shared_ptr<Item> p) : name(s), parent(p) {};

#if _DEBUG
	~Item() {cout << "destructor:" + name << endl;};
#endif
	void add(ItemPtr c) { this->children.push_back(c); };

	string myname() {
		size_t ridx = this->name.rfind(DELIMITER_CHAR);
		return this->name.substr(ridx + 1);
	};
	inline shared_ptr<Item> get_ptr(list<ItemPtr>::iterator it) {
		return  it->lock();
	};
	inline shared_ptr<Item> get_parent() {
		return  parent;
	};
	void accept(Visitor& v) {
		v.visit(*this);
	};
	const list<ItemPtr>::iterator iterator() {return children.begin();}
	const list<ItemPtr>::iterator iterator_end() {return children.end();}
};

class ItemVisitor : public Visitor {
private:
	int level = 0;
public:
	virtual void visit(Item& item) {
		for (int i = 0; i < level * 4; ++i) cout << " ";
		cout << item.myname() << endl;

		for (auto it = item.iterator(); it != item.iterator_end(); ++it) {
			shared_ptr<Item> si = item.get_ptr(it);
			level++;
			si->accept(*this);
			level--;
		}
	}
};
class LineItemVisitor : public Visitor {
private:
	string vline_last;
	string vline_not_last;
	string hline_last;
	string hline_not_last;

	inline bool is_last(Item* p, ItemPtr& last_item) {
		shared_ptr<Item> si = last_item.lock();
		return (si.get() == p);
	};
	inline bool is_higher_last(shared_ptr<Item>& self, ItemPtr& last_item) {
		shared_ptr<Item> si = last_item.lock();
		return (si.get()== self.get());
	};
public:
	LineItemVisitor(string vl,string vnl,string hl,string hnl)
		: vline_last(vl),vline_not_last(vnl),hline_last(hl),hline_not_last(hnl){};

	virtual void visit(Item& item) {

		if (item.parent != nullptr) {
			vector<string> keisen;
			keisen.push_back((is_last(&item, item.parent->children.back()))?hline_last:hline_not_last);

			auto c = item.parent;
			while (c->parent != nullptr) {
				keisen.push_back((is_higher_last(c, c->parent->children.back()))?vline_last:vline_not_last);
				c = c->parent;
			}
			for (int i = keisen.size() - 1; 0 <= i; --i ){
				cout << keisen[i];
			}
		}
		cout << item.myname() << endl;
		for (auto it = item.iterator(); it != item.iterator_end(); ++it) {
			shared_ptr<Item> si = item.get_ptr(it);
			si->accept(*this);
		}
	};
};
void create_tree_ordered(shared_ptr<Item>& top, istream& in) {

	static unordered_map< string,shared_ptr<Item> > cache;
	string full_name;

	while (true) {
		if (!std::getline(in,full_name)) break;
		if (cache[full_name] != nullptr) continue;

		size_t ridx = full_name.rfind(DELIMITER_CHAR);
		if (-1 == ridx) {
			top = cache[full_name] = shared_ptr<Item>(new Item(full_name));
		} else {
			auto parent_name = full_name.substr(0,ridx);
			auto parent = cache[parent_name];
			cache[full_name] = shared_ptr<Item>(new Item(full_name, parent));
			parent->add(cache[full_name]);
		}
	}
}
void split(vector<string>& v, const string &str, char sep)
{
	stringstream ss(str);
	string buffer;
	while( std::getline(ss, buffer, sep) ) {
		v.push_back(buffer);
	}
}
void create_tree(shared_ptr<Item>& top, istream& in) {

	static unordered_map< string,shared_ptr<Item> > cache;
	string full_name;

	while (true) {
		if (!std::getline(in,full_name)) break;

		string items = "";
		vector<string> vec;
		split(vec, full_name, DELIMITER_CHAR);

		for (auto it = vec.begin();  it != vec.end(); ++it) {

			if (items == "") {
				items = *it;
			} else {
				items = items + DELIMITER_CHAR + *it;
			}

			auto k = cache.find(items);
			if (k != cache.end()) continue;

			size_t ridx = items.rfind(DELIMITER_CHAR);

			if (-1 == ridx) {
				top = cache[items] = shared_ptr<Item>(new Item(items));
			} else {
				auto parent_name = items.substr(0,ridx);
				auto parent = cache[parent_name];
				cache[items] = shared_ptr<Item>(new Item(items, parent));
				parent->add(cache[items]);
			}
		}
	}
}
int main(int argc,char** argv) {

	int rtc = 0;
	string filename = "";
	bool line = false;
	bool mline = false;
	void (*create_tree_impl)(shared_ptr<Item>&, istream&);
	create_tree_impl = create_tree;

	for (int opt = 0,opterr = 0; (opt = getopt(argc, argv, "f:lmo")) != -1; ) {
		switch (opt) {
		case 'f':
			filename = optarg;
			break;
		case 'l':
			line = true;
			break;
		case 'm':
			mline = true;
			break;
		case 'o':
			create_tree_impl = create_tree_ordered;
			break;
		default:
			break;
		}
	}

	shared_ptr<Item> fj;
	try {
		if (filename == "") {
			create_tree_impl(fj, cin);
		} else {
			ifstream ifs(filename);
			create_tree_impl(fj,ifs);
			ifs.close();
		}
	}
	catch (exception& x) {
		cerr << x.what() << endl;
	}
	Visitor* v;
	if (line) {
		v = new LineItemVisitor("   ","|  ", "`--" ,"|--" );
	} else if (mline) {
		v = new LineItemVisitor("　", "┃","┗","┣");
	} else {
		v = new ItemVisitor();
	}
	fj->accept(*v);
	return rtc;
}
